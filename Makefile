MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -e -c

RM := rm -rf

# Automatically discover all buildable targets.
BUILDABLE_PACKAGES := $(shell \
	test -r .Makefile.packages.cache && cat .Makefile.packages.cache \
	|| \
	find . \
		-not \( -path ./Godeps -prune \) \
		-type f \
		-name "*.go" \
		-exec bash -c 'test -n "`tr "\n" " " < "{}" | grep "^package main .* func main() { .* }"`" && echo {}' \; \
	| xargs -n1 -IX bash -c '( cd "`dirname "X"`" && go list && cd - 1>/dev/null ) | sed "s|^\([^\/]*\/\)\{0,1\}\(.*\)|`dirname "X"`\/\2:`dirname "X"`|"' \
	| tee .Makefile.cache \
)

BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILD_HASH := $(shell git log -1 | head -n 1 | cut -d ' ' -f 2)


all: generate build

generate:
	set -e
	go generate ./...

build:
	set -e
	$(foreach PACKAGE_INFO,$(BUILDABLE_PACKAGES), \
		$(eval ARTIFACT = $(word 1,$(subst :, ,$(PACKAGE_INFO)))) \
		$(eval PACKAGE_PATH = $(word 2,$(subst :, ,$(PACKAGE_INFO)))) \
		godep go build -ldflags "-X main.BuildTimeStr $(BUILD_DATE) -X main.BuildCommitHash $(BUILD_HASH)" -o $(ARTIFACT) $(PACKAGE_PATH) ; \
	)


clean:
	$(RM) $(foreach PACKAGE_INFO,$(BUILDABLE_PACKAGES), \
		$(eval ARTIFACT = $(word 1,$(subst :, ,$(PACKAGE_INFO)))) \
		$(ARTIFACT) \
	)
