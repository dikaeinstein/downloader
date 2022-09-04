BINARY_NAME=

PACKAGE=main
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go env GOVERSION)

test:
	@go test -race $(TESTFLAGS) ./...

test-cover:
	@go test -coverprofile=cover.out -race $(TESTFLAGS) ./...
	@go tool cover -html=cover.out -o cover.html

lint:
	@golangci-lint run

LDFLAGS=-ldflags '-s -w \
	-X "$(PACKAGE).binaryVersion=$(VERSION)" \
	-X "$(PACKAGE).buildDate=$(BUILD_DATE)" \
	-X "$(PACKAGE).goVersion=$(GO_VERSION)" \
	-X "$(PACKAGE).gitHash=$(GIT_COMMIT_HASH)"'

build:
	@go build $(LDFLAGS) ./cmd/downloadctl

build-linux:
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) ./cmd/downloadctl

install:
	@go install $(LDFLAGS) ./cmd/downloadctl

run:
	@go run $(LDFLAGS) ./cmd/downloadctl

## send test coverage to coveralls
coveralls:
	@go run github.com/mattn/goveralls -coverprofile=cover.out -service=github

## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
