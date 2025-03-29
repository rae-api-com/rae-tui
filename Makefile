# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=rae
BINARY_UNIX=$(BINARY_NAME)_unix

# Version information
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# All
all: install

# Build
build: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./build/$(BINARY_NAME) -v ./...

install: build
	sudo install -m 755 -c ./build/rae /usr/local/bin/rae

# Clean
clean:
	$(GOCLEAN)
	rm -rf build

# Format
fmt:
	$(GOCMD) fmt ./...
	goimports -w .
	golines -w .

run:
	go run $(LDFLAGS) ./*.go

tidy:
	$(GOMOD) tidy

setup:
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: all build clean run fmt vendor deps build-linux tidy install setup
