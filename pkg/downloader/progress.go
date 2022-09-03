package downloader

import (
	"fmt"
	"math"

	"github.com/schollz/progressbar/v3"
)

// DefaultProgress just shows the progress as a percentage.
type DefaultProgress struct {
	bytesWritten       int64
	totalExpectedBytes int64
}

// Write implements the io.Writer interface.
func (dp *DefaultProgress) Write(p []byte) (int, error) {
	n := len(p)
	dp.bytesWritten += int64(n)
	percentDownloaded := float64(dp.bytesWritten) / float64(dp.totalExpectedBytes) * 100

	fmt.Printf("\rDownloading... %.0f%% complete", math.Round(percentDownloaded))
	return n, nil
}

// SetTotalBytes sets the total number of bytes to be downloaded.
func (dp *DefaultProgress) SetTotalBytes(totalBytes int64) {
	dp.totalExpectedBytes = totalBytes
}

type ProgressBar struct {
	*progressbar.ProgressBar
}

func (pb *ProgressBar) SetTotalBytes(totalBytes int64) {
	pb.ProgressBar = progressbar.DefaultBytes(totalBytes, "Downloading")
}
