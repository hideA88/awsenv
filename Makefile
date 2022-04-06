export PATH := $(CURDIR)/.bin:$(PATH)

TARGETS = awsenv

GOLANGCI_LINT = golangci-lint run
TEST = ./...
PKGNAME = $(shell go list -m)
GIT_COMMIT = $(shell git rev-parse HEAD)
VERSION =$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
BUILD=local
LDFLAGS = -ldflags "-X $(PKGNAME)/pkg/version.gitCommit=$(GIT_COMMIT) \
										-X $(PKGNAME)/pkg/version.version=$(VERSION)\
										-X $(PKGNAME)/pkg/version.build=$(BUILD)"

# command
defualt: tools help

## Install dependency and tools
setup: deps tools

## Install dependency
deps:
	go get ./...
deps.update.minor:
	go get -t -u ./...
deps.update.patch:
	go get -t -u=patch ./...
deps.tidy:
	go mod tidy

## Install tools
tools:
	export GOBIN=$(CURDIR)/.bin &&\
	go install github.com/Songmu/make2help/cmd/make2help@v0.2.0 &&\
	go install github.com/kyoh86/richgo@v0.3.10 &&\
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

## Remove build target
clean:
	rm -f $(TARGETS)
	rm -rf dist
	rm -rf tmp

## Build app
build: clean deps
	go build $(LDFLAGS) -o dist/$(TARGETS) main.go

## Check code format
check:
	$(GOLANGCI_LINT) ./...

## Fix code
fix:
	$(GOLANGCI_LINT) --fix ./...

## Run test
test: tools
	mkdir -p tmp
	richgo test -race -coverprofile=tmp/coverage.txt -covermode=atomic $(TEST)

## Show help
help:
	@make2help $(MAKEFILE_LIST)

NO_PHONY = /^:/
PHONY := $(shell cat $(MAKEFILE_LIST) | awk -F':' '/^[a-z0-9_.-]+:/ && !$(NO_PHONY) {print $$1}')
.PHONY: $(PHONY)

show_phony:
	@echo $(PHONY)
