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
- Pure domain packages go under `internal/` (none exist yet). `cmd/genesis`
  stays a thin CLI over them: argument parsing, output and exit codes only.

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

## Binaries

- Build scratch binaries outside `$PATH`:
  `go build -o $TMPDIR/genesis-x ./cmd/genesis && $TMPDIR/genesis-x version`.
  `~/bin/genesis` is one slot shared by every worktree, so it may not be yours.
- Never `cp` a Go binary on arm64. Its ad-hoc signature is bound to the path
  it was linked at; a copy gets SIGKILLed (exit 137, no output). Recover one
  with `codesign --force --sign - <path>`.
- `make install` builds in place into `PREFIX` for that reason; keep it that
  way.
