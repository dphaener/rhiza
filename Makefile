BINARY  := genesis
PREFIX  ?= $(HOME)/bin
VERSION := $(shell git describe --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

GOFILES := $(shell find . -name '*.go' -not -path './vendor/*' -not -path './.shipit/*')

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build ./genesis in the repo root
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/genesis

.PHONY: install
install: ## Build straight into $(PREFIX) (default ~/bin)
# Built in place at the destination, never built-then-copied. The Go linker
# gives the binary an ad-hoc signature bound to the path it was written to, so
# on arm64 a `cp` elsewhere invalidates it and the kernel SIGKILLs the copy
# (exit 137, no output). Recover a copied binary with:
#   codesign --force --sign - $(PREFIX)/$(BINARY)
	@mkdir -p $(PREFIX)
	go build -ldflags "$(LDFLAGS)" -o $(PREFIX)/$(BINARY) ./cmd/genesis

.PHONY: test
test: ## Run the full test suite
	go test ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format all Go files in place
	gofmt -w $(GOFILES)

.PHONY: fmt-check
fmt-check: ## Fail if any Go file is unformatted
# `gofmt -l` exits 0 even when it lists files, so a non-empty listing has to
# be turned into a failure by hand.
	@unformatted=$$(gofmt -l $(GOFILES)); \
	if [ -n "$$unformatted" ]; then \
		echo "unformatted files:"; echo "$$unformatted"; exit 1; \
	fi

.PHONY: check
check: fmt-check vet test ## Run everything CI runs
# The single entry point .github/workflows/ci.yml runs, so local and CI
# runs cannot drift apart.

.PHONY: clean
clean: ## Remove the locally built binary
	rm -f $(BINARY)

.PHONY: help
help: ## List targets
	@grep -hE '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk -F':.*?## ' '{printf "  %-12s %s\n", $$1, $$2}'
