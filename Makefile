# Default Go toolchain
GO              ?= go
GOLANGCI_LINT   ?= golangci-lint
PKG             := github.com/tarantool/go-tlog

.PHONY: all test test-race test-coverage lint fmt tidy examples help

all: test

## Run tests
test:
	$(GO) test ./...

## Run tests with race detector
test-race:
	$(GO) test -race ./...

## Run tests with coverage
test-coverage:
	$(GO) test -covermode=atomic -coverprofile=coverage.out ./...

## Run golangci-lint
lint:
	$(GOLANGCI_LINT) run ./...

## Format source code
fmt:
	$(GO) fmt ./...

## Tidy go.mod / go.sum
tidy:
	$(GO) mod tidy

## Run all _examples to ensure they compile and run without panic
examples:
	$(GO) run ./_examples/stdout
	$(GO) run ./_examples/stderr >/dev/null 2>&1 || true
	$(GO) run ./_examples/file
	$(GO) run ./_examples/multi

## Show available targets
help:
	@echo "Available targets:"
	@echo "  make test           - run tests"
	@echo "  make test-race      - run tests with -race"
	@echo "  make test-coverage  - run tests with coverage"
	@echo "  make lint           - run golangci-lint"
	@echo "  make fmt            - format sources (gofmt)"
	@echo "  make tidy           - go mod tidy"
	@echo "  make examples       - run all examples"
