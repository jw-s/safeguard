.PHONY: all clean build
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES=$(shell go list ./... | grep -v /vendor/)
export GO111MODULE=on

all: build

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o ./bin/safeguard ./cmd
clean:
	rm -rf bin/*

fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

lint:
	@golint ./pkg/...

vet:
	@go vet ./pkg/...
