# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Ignore the `~/code/.claude/rules` files

The Rails/Ruby rules in `~/code/.claude/rules/` (`stack.md`, `global.md`,
`ui-standards.md`, `boundaries.md`) are loaded only because this repo sits
under `~/code`. They do not apply to Genesis. Ignore them: there is no Rails,
Ruby, RSpec, Tailwind or ERB here.

## Stack and shape

- Go, one binary: `genesis`. The future daemon is the same binary, not a
  second one.
- Pure domain packages go under `internal/`. `cmd/genesis` stays a thin CLI
  over them: argument parsing, output and exit codes only.
- Operations live in `internal/api`: transport-independent, one registry,
  schemas reflected from the Go request/response types. Its interim model
  (over `changeset.FoldLog`) is to be replaced by projection.
- `cmd/genesis` is a generic frame: it maps argv to requests through the
  request structs' tags (`arg`, `flag`, `desc`). `genesis <cmd> --json` and
  `genesis api call <op>` reach the same handler and print the same bytes.

## Commands

```bash
make            # build ./genesis in the repo root
make check      # everything CI runs (fmt-check + vet + test)
make test       # go test ./...
make fmt        # gofmt every Go file in place
make install    # build into ~/bin (PREFIX=... to override)
```

Run `make check` before calling any change done.

## Tests

- Scenarios are testscript files in `cmd/genesis/testdata/script/*.txtar`.
  They run in-process via `testscript.Main` in `cmd/genesis/main_test.go`,
  never against the `genesis` on `$PATH`.
- Git inside scripts is hermetic: the `Setup` in `cmd/genesis/main_test.go`
  gives each script its own `HOME` and global git config.
- No test touches the network, Rex, agents or `gh`.
- Nothing binds a unix socket in the default test run.
- Every script runs with `GENESIS_TEST_NOW` and `GENESIS_TEST_SEED` set, so
  IDs and times are deterministic. Both are offset by the entries in the log
  (plus one for `repo.json`), so each run draws fresh IDs and a `--dry-run`
  predicts the next real run's.
- The `fold-ok` script command checks that `.genesis` reads and folds with
  no problems.
- `describe-golden.txtar` holds `describe --json`. Regenerate it when a
  request or response type changes:
  `go test ./cmd/genesis -run TestScripts/describe-golden -update`.

## Binaries

- Build scratch binaries outside `$PATH`:
  `go build -o $TMPDIR/genesis-x ./cmd/genesis && $TMPDIR/genesis-x version`.
  `~/bin/genesis` is one slot shared by every worktree, so it may not be yours.
- Never `cp` a Go binary on arm64. Its ad-hoc signature is bound to the path
  it was linked at; a copy gets SIGKILLed (exit 137, no output). Recover one
  with `codesign --force --sign - <path>`.
- `make install` builds in place into `PREFIX` for that reason; keep it that
  way.
