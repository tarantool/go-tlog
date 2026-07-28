# Default Go toolchain
GO              ?= go
GOLANGCI_LINT   ?= golangci-lint
PKG             := github.com/tarantool/go-tlog

# Pinned so that a local run and CI agree. Keep in sync with the
# GOLANGCI_LINT_VERSION env in .github/workflows/testing.yml.
GOLANGCI_LINT_VERSION ?= v2.11.4

COVERAGE_PROFILE ?= coverage.out
# internal/slog is a vendored fork of the standard library's log/slog. It is
# covered by the Go project's own tests; measuring it here would only dilute
# the number for code this repository actually owns.
COVERAGE_EXCLUDE := $(PKG)/internal/slog/

.PHONY: all test test-race test-coverage lint fmt spell vulncheck generate-check \
	install-lint tidy help

all: test

## Run tests.
test:
	$(GO) test ./...

## Run tests with race detector.
test-race:
	$(GO) test -race ./...

## Run tests with coverage.
test-coverage:
	$(GO) test -covermode=atomic -coverprofile=$(COVERAGE_PROFILE).tmp ./...
	grep -v '$(COVERAGE_EXCLUDE)' $(COVERAGE_PROFILE).tmp > $(COVERAGE_PROFILE)
	rm -f $(COVERAGE_PROFILE).tmp
	$(GO) tool cover -func=$(COVERAGE_PROFILE) | tail -1

## Run golangci-lint.
lint:
	$(GOLANGCI_LINT) run ./... --config=.golangci.yaml

## Format source code.
fmt:
	$(GOLANGCI_LINT) fmt --config=.golangci.yaml

## Check spelling in code and documentation.
spell:
	codespell

## Check dependencies for known vulnerabilities.
vulncheck:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Fail if `go generate` leaves a diff.
generate-check:
	$(GO) generate ./...
	git diff --exit-code

## Install the pinned golangci-lint.
install-lint:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## Tidy go.mod / go.sum.
tidy:
	$(GO) mod tidy

## Show available targets.
help:
	@echo "Available targets:"
	@echo "  make test            - run tests"
	@echo "  make test-race       - run tests with -race"
	@echo "  make test-coverage   - run tests with coverage"
	@echo "  make lint            - run golangci-lint"
	@echo "  make fmt             - format sources"
	@echo "  make spell           - run codespell"
	@echo "  make vulncheck       - run govulncheck"
	@echo "  make generate-check  - verify go generate leaves no diff"
	@echo "  make install-lint    - install the pinned golangci-lint"
	@echo "  make tidy            - go mod tidy"
