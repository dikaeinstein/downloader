# Downloader

Downloads files using the given URL. It supports parallel file download i.e it can download a file in chunks and then merge it.

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
