# Downloader

Downloads files using the given URL. It supports sync and parallel file download i.e it can download a file in chunks and then merge it.
It provides a library that can be used in other projects and also a cli tool (downloadctl).

[![CI/CD](https://github.com/dikaeinstein/downloader/actions/workflows/ci-cd.yaml/badge.svg)](https://github.com/dikaeinstein/downloader/actions/workflows/ci-cd.yaml)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/downloader/badge.svg?branch=main)](https://coveralls.io/github/dikaeinstein/downloader?branch=main)

## Download and Install

To download and install a specific version of godl, copy and paste the installation command:

```bash
curl -s https://raw.githubusercontent.com/dikaeinstein/downloader/master/get.sh | sh -s -- v0.1.0
```

## Install with Go 1.18+

The binary can be installed with `go install`. The binary is placed in $GOPATH/bin, or in $GOBIN if set:

```bash
go install github.com/dikaeinstein/downloader/cmd
```

## Build From Source

Prerequisites for building from source are:

- make
- Go 1.18+

```bash
git clone https://github.com/dikaeinstein/downloader
cd downloader
make build
```
