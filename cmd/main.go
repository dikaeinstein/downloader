package main

import (
	"github.com/dikaeinstein/downloader/internal/cli"
)

// injected as ldflags during go build
var (
	buildDate     = "unknown"
	gitHash       = "none"
	binaryVersion = "v0.0.0"
	goVersion     = "unknown"
)

func main() {
	opt := cli.Option{}

	opt.Version.BuildDate = buildDate
	opt.Version.GitHash = gitHash
	opt.Version.BinaryVersion = binaryVersion
	opt.Version.GoVersion = goVersion

	cli.Run(opt)
}
