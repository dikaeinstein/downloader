package downloader_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/downloader/internal/pkg/fsys"
	"github.com/dikaeinstein/downloader/pkg/downloader"
	"github.com/dikaeinstein/downloader/pkg/hash"
)

func TestDownloader(t *testing.T) {
	testCases := []struct {
		desc               string
		url                string
		contentLength      int64
		filename           string
		downloadedFilename string
		resBody            *bytes.Buffer
		statusCode         int
	}{
		{
			desc:               "can download remote file",
			url:                "",
			contentLength:      17,
			filename:           "testFile",
			downloadedFilename: "testFile",
			resBody:            bytes.NewBufferString("This is test data"),
			statusCode:         http.StatusOK,
		},
		{
			desc:               "can extract filename from url",
			url:                "https://example.com/testFile",
			contentLength:      17,
			filename:           "",
			downloadedFilename: "testFile",
			resBody:            bytes.NewBufferString("This is test data"),
			statusCode:         http.StatusOK,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			resContent := tC.resBody.Bytes()
			fakeRoundTripper := roundTripFunc(
				func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode:    http.StatusOK,
						Body:          io.NopCloser(tC.resBody),
						ContentLength: int64(tC.resBody.Len()),
					}
				},
			)

			testClient := &http.Client{Transport: fakeRoundTripper}
			inMemFS := fsys.NewInMemFS(make(fstest.MapFS))
			hasher := hash.FakeHasher{}
			verifier := fakeHashVerifier{}
			progresser := downloader.DefaultProgress{}

			d, err := downloader.New(
				".", testClient, inMemFS, hasher, progresser, verifier,
			)
			require.NoError(t, err)

			err = d.Download(
				context.Background(), tC.url, tC.filename, "test.256")
			require.NoError(t, err)

			require.Equal(t, exists(t, inMemFS, tC.downloadedFilename), true)

			content := inMemFS.Content(tC.downloadedFilename)
			require.Equal(t, tC.contentLength, int64(content.Len()))
			require.Equal(t, string(resContent), content.String())
		})
	}
}

// The RoundTripFunc type is an adapter to allow the use of
// ordinary functions as  net/http.RoundTripper. If f is a function
// with the appropriate signature, RoundTripFunc(f) is a
// RoundTripper that calls f.
type roundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type fakeHashVerifier struct{}

func (fakeHashVerifier) Verify(input io.Reader, hex string) error {
	return nil
}

func exists(t *testing.T, fsys fs.StatFS, filename string) bool {
	t.Helper()

	_, err := fsys.Stat(filename)
	return !errors.Is(err, fs.ErrNotExist)
}
