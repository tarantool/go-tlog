# Default Go toolchain
GO              ?= go
GOLANGCI_LINT   ?= golangci-lint
PKG             := github.com/tarantool/go-tlog

.PHONY: all test test-race test-coverage lint fmt tidy help

all: test

## Run tests.
test:
	$(GO) test ./...

## Run tests with race detector.
test-race:
	$(GO) test -race ./...

## Run tests with coverage.
test-coverage:
	$(GO) test -covermode=atomic -coverprofile=coverage.out ./...

## Run golangci-lint.
lint:
	$(GOLANGCI_LINT) run ./... --config=.golangci.yaml

## Format source code.
fmt:
	$(GO) fmt ./...

## Tidy go.mod / go.sum.
tidy:
	$(GO) mod tidy

## Show available targets.
help:
	@echo "Available targets:"
	@echo "  make test           - run tests"
	@echo "  make test-race      - run tests with -race"
	@echo "  make test-coverage  - run tests with coverage"
	@echo "  make lint           - run golangci-lint"
	@echo "  make fmt            - format sources (gofmt)"
	@echo "  make tidy           - go mod tidy"
