BINARY  := rhiza
PREFIX  ?= $(HOME)/bin
VERSION := $(shell git describe --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

GOFILES := $(shell find . -name '*.go' -not -path './vendor/*' -not -path './.shipit/*')

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build ./rhiza in the repo root
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/rhiza

.PHONY: install
install: ## Build straight into $(PREFIX) (default ~/bin)
# Built in place at the destination, never built-then-copied. The Go linker
# gives the binary an ad-hoc signature bound to the path it was written to, so
# on arm64 a `cp` elsewhere invalidates it and the kernel SIGKILLs the copy
# (exit 137, no output). Recover a copied binary with:
#   codesign --force --sign - $(PREFIX)/$(BINARY)
	@mkdir -p $(PREFIX)
	go build -ldflags "$(LDFLAGS)" -o $(PREFIX)/$(BINARY) ./cmd/rhiza

.PHONY: update
update: install ## Rebuild into $(PREFIX), then rewrite the Rex script and reload Rex
# Runs the binary install just built, in place, so the Rex script always comes
# from the same build as the binary. Nothing is copied (arm64 signature, see
# install).
	$(PREFIX)/$(BINARY) rex install

.PHONY: test
# The timing budgets measure reads against fixed limits, so they run on
# their own after the rest of the suite: run beside cmd/rhiza's scripts they
# measured their neighbours (PR #210 failed twice, then passed, on one SHA).
BUDGET_TESTS = ^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestGeneratedModel)$$
BUDGET_PKGS = ./internal/api ./internal/index ./internal/slice

# -trimpath keeps each worktree's absolute path out of the compiled
# packages and test binaries, so a second worktree at the same code reuses
# the first one's builds, and a package whose tests read no file in the repo
# (and not $$PWD) is a test-cache hit there too. Tests that do read repo
# files (testdata, the corpus, rapid's testdata/rapid lookup) still miss: go
# keys those results on the files' absolute paths and mtimes.
TEST_FLAGS = -trimpath -timeout 20m

# cmd/rhiza's scripts run in parallel and keep every core busy, so on one
# machine only one worktree runs them at a time: concurrent `make test`s
# (one per Rhiza run worktree) queue on this lock instead of thrashing. The
# rest of the suite mostly waits on git, fsync and sleeps rather than CPU,
# so it runs unlocked, beside the scripts. The timing budgets take the lock
# too, so they never run beside another suite's scripts. flock(1) is
# util-linux; where it is missing (macOS) nothing is locked. TEST_LOCK=
# skips the lock.
TEST_LOCK ?= $(dir $(shell go env GOCACHE))rhiza-test.lock

test: ## Run the full test suite, then the timing budgets alone, and print the read-budget timings
# TestReadBudget appends its timings to the file named by RHIZA_BUDGET_LOG.
# They print after the run whether the tests passed or failed, so CI logs
# (and local runs) always show the numbers. go test's exit status is make's.
# The packages other than cmd/rhiza run in the background, their output
# held in a file and printed once cmd/rhiza is done, so the two never
# interleave. The shell holds the lock on fd 9 (lock, unlock); every go test
# gets 9>&-, so no process a test leaves behind can keep holding it.
	@lockfile='$(TEST_LOCK)'; \
	lock() { \
		if [ -n "$$lockfile" ] && command -v flock >/dev/null 2>&1 && (: >>"$$lockfile") 2>/dev/null; then \
			exec 9>>"$$lockfile"; \
			if ! flock -n 9; then \
				echo "Waiting for another make test on this machine before running $$1 (lock: $$lockfile)"; \
				start=$$(date +%s); flock 9; \
				echo "Got the test lock after $$(( $$(date +%s) - start ))s"; \
			fi; \
		fi; \
	}; \
	unlock() { exec 9>&-; }; \
	log=$$(mktemp); rest=$$(mktemp); \
	go test $(TEST_FLAGS) -skip '$(BUDGET_TESTS)' $$(go list ./... | grep -v '/cmd/rhiza$$') >$$rest 2>&1 9>&- & \
	pid=$$!; \
	lock "the cmd/rhiza scripts"; \
	go test $(TEST_FLAGS) -skip '$(BUDGET_TESTS)' ./cmd/rhiza 9>&-; \
	status=$$?; \
	unlock; \
	wait $$pid || status=1; \
	cat $$rest; \
	echo "Timing budgets, run alone after the suite"; \
	lock "the timing budgets"; \
	RHIZA_BUDGET_LOG=$$log go test $(TEST_FLAGS) -count=1 -p 1 -run '$(BUDGET_TESTS)' $(BUDGET_PKGS) 9>&- || status=1; \
	unlock; \
	echo "Read-budget timings"; \
	if [ -s "$$log" ]; then cat "$$log"; else echo "(no timings: $$log is missing or empty)"; fi; \
	rm -f $$log $$rest; exit $$status

# One target per CI job, so the jobs run in parallel and each is a slice of
# what `make test` runs. Together test-scripts over every shard, test-api
# over every shard, test-rest and test-budgets cover the same tests as
# `make test`.
#
# SHARD is "i/n" (1-based). GO_TEST is the go test command; the tests
# replace it to see the arguments a target would pass.
SHARD ?= 1/1
GO_TEST ?= go test
SKIP_BUDGETS = $(BUDGET_TESTS)

# LOCK_SH defines lock and unlock, as in test: the machine-wide lock on fd 9.
LOCK_SH = lockfile='$(TEST_LOCK)'; \
	lock() { \
		if [ -n "$$lockfile" ] && command -v flock >/dev/null 2>&1 && (: >>"$$lockfile") 2>/dev/null; then \
			exec 9>>"$$lockfile"; \
			if ! flock -n 9; then \
				echo "Waiting for another make test on this machine before running $$1 (lock: $$lockfile)"; \
				start=$$(date +%s); flock 9; \
				echo "Got the test lock after $$(( $$(date +%s) - start ))s"; \
			fi; \
		fi; \
	}; \
	unlock() { exec 9>&-; }

.PHONY: test-scripts
test-scripts: ## Run shard SHARD=i/n of cmd/rhiza's scripts, timing budgets skipped
	@$(LOCK_SH); \
	lock "the cmd/rhiza scripts"; \
	RHIZA_TEST_SHARD='$(SHARD)' $(GO_TEST) $(TEST_FLAGS) -run '^TestScripts$$' -skip '$(SKIP_BUDGETS)' ./cmd/rhiza 9>&-; \
	status=$$?; unlock; exit $$status

.PHONY: test-api
test-api: ## Run shard SHARD=i/n of internal/api's top-level tests, timing budgets excluded
# The tests come from go test -list, budgets dropped, and the shard takes
# every n-th of them (round-robin) as an anchored -run regex.
	@i='$(SHARD)'; n=$${i#*/}; i=$${i%/*}; \
	case "$$i$$n" in *[!0-9]*|'') echo "SHARD=$(SHARD): want i/n with 1 <= i <= n" >&2; exit 1;; esac; \
	if [ "$$i" -lt 1 ] || [ "$$i" -gt "$$n" ]; then echo "SHARD=$(SHARD): want i/n with 1 <= i <= n" >&2; exit 1; fi; \
	names=$$(go test -trimpath -list '^Test' ./internal/api | grep '^Test' | grep -Ev '$(SKIP_BUDGETS)' | LC_ALL=C sort) || exit 1; \
	re=$$(printf '%s\n' "$$names" | awk -v i="$$i" -v n="$$n" 'NR % n == i % n { printf "%s%s", sep, $$0; sep = "|" }'); \
	if [ -z "$$re" ]; then echo "shard $(SHARD) of internal/api has no tests"; exit 0; fi; \
	$(GO_TEST) $(TEST_FLAGS) -run "^($$re)\$$" ./internal/api

.PHONY: test-rest
test-rest: ## Run every package but internal/api, timing budgets skipped
	$(GO_TEST) $(TEST_FLAGS) -skip '$(BUDGET_TESTS)|^TestScripts$$' $$(go list ./... | grep -v '/internal/api$$')

.PHONY: test-budgets
test-budgets: ## Run the timing budgets alone, and print the read-budget timings
	@log=$$(mktemp); $(LOCK_SH); \
	lock "the timing budgets"; \
	RHIZA_BUDGET_LOG=$$log $(GO_TEST) $(TEST_FLAGS) -count=1 -p 1 -run '$(SKIP_BUDGETS)' $(BUDGET_PKGS) 9>&-; \
	status=$$?; unlock; \
	echo "Read-budget timings"; \
	if [ -s "$$log" ]; then cat "$$log"; else echo "(no timings: $$log is missing or empty)"; fi; \
	rm -f $$log; exit $$status

.PHONY: lint
lint: fmt-check vet self-check ## Run fmt-check, vet and self-check

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
check: fmt-check vet self-check test ## Run everything CI runs
# The single entry point .github/workflows/ci.yml runs, so local and CI
# runs cannot drift apart.

.PHONY: self-check
self-check: ## Run rhiza check on this repo with this commit's build
# Rhiza keeps its own .rhiza/, so the build under test checks it: a stale
# rendered/ or a broken log fails make. The binary is built in place in a
# fresh temp folder (never copied there, for the arm64 signature reason under
# install), so it never touches ./rhiza or ~/bin. The folder is removed
# whatever check returns, and check's exit status is make's.
	@dir=$$(mktemp -d); \
	go build -ldflags "$(LDFLAGS)" -o $$dir/$(BINARY) ./cmd/rhiza \
		&& $$dir/$(BINARY) check; \
	status=$$?; rm -rf $$dir; exit $$status

.PHONY: docs
docs: ## Regenerate docs/reference with this commit's build
# The project's generate command (rhiza.ship.generate): build and run rhiza
# docs the way self-check runs check, so the reference always comes from the
# branch's own code and never from a stale binary on $$PATH.
	@dir=$$(mktemp -d); \
	go build -ldflags "$(LDFLAGS)" -o $$dir/$(BINARY) ./cmd/rhiza \
		&& $$dir/$(BINARY) docs; \
	status=$$?; rm -rf $$dir; exit $$status

.PHONY: clean
clean: ## Remove the locally built binary
	rm -f $(BINARY)

.PHONY: help
help: ## List targets
	@grep -hE '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk -F':.*?## ' '{printf "  %-12s %s\n", $$1, $$2}'
