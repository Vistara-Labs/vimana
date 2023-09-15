# Basic Makefile for Go CLI project

# Name of your project (binary)
BINARY_NAME=vimana

# Version and build time information
VERSION=0.0.1
BUILD_TIME=$(shell date +%FT%T%z)

# Flags passed to `go build`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default command when simply typing `make`
all: build

# Compiles the binary
build:
	@echo "Building..."
	go build ${LDFLAGS} -o ${BINARY_NAME} -v ./...
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
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 -v
	mkdir -p ${HOME}/.vimana
	cp config.toml ${HOME}/.vimana/config.toml

.PHONY: all build clean test install build-linux
