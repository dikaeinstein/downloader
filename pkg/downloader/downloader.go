package downloader

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dikaeinstein/downloader/internal/pkg/fsys"
)

// Hasher is an interface to retrieve a file's checksum.
type Hasher interface {
	// Hash returns the checksum at the given path.
	// path can be a local file containing the checksum
	// or a remote URL to download the checksum from.
	Hash(ctx context.Context, path string) (string, error)
}

// HashVerifier is an interface to verify a file's checksum.
type HashVerifier interface {
	// Verify ensure the checksum of the given io.Reader
	//  matches the given hash.
	Verify(in io.Reader, want string) error
}

// Progresser is an interface to report the progress of a download.
type Progresser interface {
	// Progress reports the progress of a download.
	Progress(percentDownloaded float64)
}

// Downloader represents a downloader that can download files from a URL,
// and verify the checksum of the downloaded file.
type Downloader struct {
	downloadDir string
	fsys        fs.FS
	httpClient  *http.Client
	hasher      Hasher
	progresser  Progresser
	verifier    HashVerifier
}

// New returns a new Downloader.
func New(
	dlDir string,
	httpClient *http.Client,
	fsys fs.FS,
	hasher Hasher,
	progresser Progresser,
	verifier HashVerifier,
) (*Downloader, error) {
	err := os.MkdirAll(dlDir, os.ModeDir)
	if err != nil {
		return nil, err
	}

	dl := &Downloader{
		downloadDir: dlDir,
		fsys:        fsys,
		httpClient:  httpClient,
		hasher:      hasher,
		progresser:  progresser,
		verifier:    verifier,
	}

	return dl, nil
}

// Download downloads a file from a URL,
// and verifies the checksum of the downloaded file.
func (dl *Downloader) Download(
	ctx context.Context, url, filename, hashPath string,
) error {
	return dl.syncDownload(ctx, url, filename, hashPath)
}

// syncDownload downloads a file from a URL synchronously.
func (dl *Downloader) syncDownload(
	ctx context.Context, url, filename, hashPath string,
) error {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, http.NoBody,
	)
	if err != nil {
		return err
	}

	res, err := dl.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		return fmt.Errorf(
			"download failed: %s: %s: %s", url, res.Status, string(data))
	}

	if filename == "" {
		filename = filepath.Base(url)
	}
	downloadPath := filepath.Join(dl.downloadDir, filename)

	// Create the file with tmp extension. So we don't overwrite until
	// the file is completely downloaded.
	tmp, err := fsys.Create(dl.fsys, downloadPath+".tmp")
	if err != nil {
		return err
	}
	defer tmp.Close()

	tmpFile, ok := tmp.(io.Writer)
	if !ok {
		return fmt.Errorf("invalid writer: %T", tmp)
	}

	// Create the bytesCountWriter used for counting response bytes
	bcw := &bytesCountWriter{
		totalExpectedBytes: res.ContentLength,
		progresser:         dl.progresser,
	}

	n, err := io.Copy(tmpFile, io.TeeReader(res.Body, bcw))
	if err != nil {
		return err
	}

	if res.ContentLength != -1 && res.ContentLength != n {
		return fmt.Errorf("copied %v bytes; expected %v", n, res.ContentLength)
	}

	wantHex, err := dl.hasher.Hash(ctx, hashPath)
	if err != nil {
		return err
	}

	fmt.Println("\nverifying checksum")
	f, err := dl.fsys.Open(downloadPath + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	err = dl.verifier.Verify(f, wantHex)
	if err != nil {
		return fmt.Errorf("error verifying checksum of %v: %v", tmpFile, err)
	}

	fmt.Println("checksums matched!")

	// Rename the temporary file once fully downloaded
	return fsys.Rename(dl.fsys, downloadPath+".tmp", downloadPath)
}

// countWriter counts the number of bytes written to it.
type bytesCountWriter struct {
	bytesWritten       int64
	totalExpectedBytes int64
	progresser         Progresser
}

func (bwc *bytesCountWriter) Write(p []byte) (int, error) {
	n := len(p)
	bwc.bytesWritten += int64(n)
	percentDownloaded := float64(bwc.bytesWritten) / float64(bwc.totalExpectedBytes) * 100

	bwc.progresser.Progress(percentDownloaded)
	return n, nil
}
