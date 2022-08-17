# Downloader

Downloads files using the given URL. It supports parallel file download i.e it can download a file in chunks and then merge it.
It provides a library that can be used in other projects and also a cli tool (downloaderctl).

[![CI/CD](https://github.com/dikaeinstein/downloader/actions/workflows/ci-cd.yaml/badge.svg)](https://github.com/dikaeinstein/downloader/actions/workflows/ci-cd.yaml)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/downloader/badge.svg?branch=main)](https://coveralls.io/github/dikaeinstein/downloader?branch=main)

## Features

- Synchronous file download
- Parallel file download
- Auto select file download method based on file size
<!-- - Download multiple files (either sync or parallel) -->

## Build From Source

Prerequisites for building from source are:

- make
- Go 1.18+

```bash
git clone https://github.com/dikaeinstein/downloader
cd downloader
make build
```
