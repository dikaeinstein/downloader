BINARY_NAME=downloadctl

PACKAGE=main
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go env GOVERSION)

test:
	@go test -race $(TESTFLAGS) ./...

test-cover:
	@go test -coverprofile=cover.out -race $(TESTFLAGS) ./...

lint:
	@golangci-lint run

LDFLAGS=-ldflags '-s -w \
	-X "$(PACKAGE).binaryVersion=$(VERSION)" \
	-X "$(PACKAGE).buildDate=$(BUILD_DATE)" \
	-X "$(PACKAGE).goVersion=$(GO_VERSION)" \
	-X "$(PACKAGE).gitHash=$(GIT_COMMIT_HASH)"'

build:
	@go build -a $(LDFLAGS) -o $(BINARY_NAME) cmd/main.go

run:
	@go run -a $(LDFLAGS) cmd/main.go

## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi

