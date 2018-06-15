.PHONY: all clean build

all: build

build:
	GOOS=linux GOARCH=amd64 go build -o bin/safeguard ./cmd

clean:
	rm -rf bin/*