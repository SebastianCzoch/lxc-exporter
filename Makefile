BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILD_HASH := $(shell git log -1 | head -n 1 | cut -d ' ' -f 2)

all: generate build

generate:
	set -e
	go get github.com/tools/godep
	go generate ./...

build:
	set -e
	$(GOPATH)/bin/godep go build -ldflags "-X main.BuildTimeStr=$(BUILD_DATE) -X main.BuildCommitHash=$(BUILD_HASH)"
