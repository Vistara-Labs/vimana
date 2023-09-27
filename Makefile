# Basic Makefile for Go CLI project

# Name of your project (binary)
BINARY_NAME=vimana

# Version and build time information
VERSION=0.0.1
BUILD_TIME=$(shell date +%FT%T%z)

# Flags passed to `go build`
LDFLAGS=-ldflags "-X vimana.Version=${VERSION} -X vimana.BuildTime=${BUILD_TIME}"

# Default command when simply typing `make`
all: build

# Compiles the binary for mac
build:
	@echo "Building for Mac..."
	mkdir -p ${BINARY_NAME}-darwin-arm64
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64 -v ./...
	mkdir -p ${HOME}/.vimana
	cp config.toml ${HOME}/.vimana/config.toml

# Cleans our project: deletes binaries
clean:
	@echo "Cleaning..."
	go clean
	rm -f ${BINARY_NAME}

# Runs tests
test:
	@echo "Testing..."
	go test -v ./...

# Installs our project: copies binaries to $GOPATH/bin
install:
	@echo "Installing..."
	go install ${LDFLAGS} ./...
	mkdir -p ${HOME}/.vimana
	cp config.toml ${HOME}/.vimana/config.toml

# Cross compilation
build-linux:
	@echo "Building for Linux..."
	mkdir -p ${BINARY_NAME}-linux-amd64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 -v ./...
	mkdir -p ${BINARY_NAME}-linux-arm64
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-arm64 -v ./...
	mkdir -p ${HOME}/.vimana
	cp config.toml ${HOME}/.vimana/config.toml

.PHONY: all build clean test install build-linux
# tar -czvf vimana_bins21.tar.gz vimana-linux-amd64/