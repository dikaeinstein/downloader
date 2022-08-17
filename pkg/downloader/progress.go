package downloader

import (
	"fmt"
	"math"
)

type DefaultProgress struct{}

func (DefaultProgress) Progress(percentDownloaded float64) {
	fmt.Printf("\rDownloading... %.0f%% complete", math.Round(percentDownloaded))
}
