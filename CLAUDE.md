# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Ignore the `~/code/.claude/rules` files

The Rails/Ruby rules in `~/code/.claude/rules/` (`stack.md`, `global.md`,
`ui-standards.md`, `boundaries.md`) are loaded only because this repo sits
under `~/code`. They do not apply to Rhiza. Ignore them: there is no Rails,
Ruby, RSpec, Tailwind or ERB here.

## Stack and shape

- Go, one binary: `rhiza`. The future daemon is the same binary, not a
  second one.
- Pure domain packages go under `internal/`. `cmd/rhiza` stays a thin CLI
  over them: argument parsing, output and exit codes only.
- Operations live in `internal/api`: transport-independent, one registry,
  schemas reflected from the Go request/response types. Its interim model
  (over `changeset.FoldLog`) is to be replaced by projection.
- `cmd/rhiza` is a generic frame: it maps argv to requests through the
  request structs' tags (`arg`, `flag`, `desc`). `rhiza <cmd> --json` and
  `rhiza api call <op>` reach the same handler and print the same bytes.

## Commands

```bash
make            # build ./rhiza in the repo root
make check      # everything CI runs (fmt-check + vet + test + self-check)
make self-check # rhiza check on this repo, with this commit's build
make test       # go test ./...
make fmt        # gofmt every Go file in place
make install    # build into ~/bin (PREFIX=... to override)
make update     # make install, then rewrite the Rex script and reload Rex
```

Agents in runs test only what their change touches, and never run
`make check`, `make test` or a bare `go test ./...`: the PR's GitHub CI runs
`make check` before every merge. The scoped commands are:

- `go test -trimpath ./internal/<pkg>/...` for each package you changed;
- `go test -trimpath ./cmd/rhiza -run TestScripts/<name>` for each affected
  testscript;
- `make self-check` for model or `.rhiza/` changes.

While CI is unstable (main's CI is not reliably green, or a red is not
explained by your change: a flaky test, an infra error, a red that passes on
rerun), do not push or rerun CI to see. Reproduce the failing test locally,
repeatedly if it is a flake, fix it or quarantine it with a bug filed, and
push only after it and the touched packages pass locally.

Outside a run (a person at a terminal), run `make check` before calling a
change done.

Outside a run, run the full suite through `make test` (or `make check`), not a bare
`go test ./...`: it builds with `-trimpath`, so packages another worktree
already tested at the same code are cache hits, and it runs cmd/rhiza's
scripts (and the timing budgets) under a machine-wide lock (`flock`,
Linux), so parallel runs' suites queue for the CPU-heavy part instead of
thrashing. "Waiting for another make test" means yours is queued. Scope
ad-hoc runs to the package you changed.

Rhiza keeps its own `.rhiza/` here, and `make check` runs `rhiza check`
on it. After landing a change set, run `rhiza render` and commit
`.rhiza/rendered/`, or `make check` fails on stale rendered files.

## Tests

- Scenarios are testscript files in `cmd/rhiza/testdata/script/*.txtar`.
  They run in-process via `testscript.Main` in `cmd/rhiza/main_test.go`,
  never against the `rhiza` on `$PATH`.
- Git inside scripts is hermetic: the `Setup` in `cmd/rhiza/main_test.go`
  gives each script its own `HOME` and global git config.
- No test touches the network, Rex, agents or `gh`.
- Nothing binds a unix socket in the default test run.
- Every script runs with `RHIZA_TEST_NOW` and `RHIZA_TEST_SEED` set, so
  IDs and times are deterministic. Both are offset by the entries in the log
  (plus one for `repo.json`), so each run draws fresh IDs and a `--dry-run`
  predicts the next real run's.
- The `fold-ok` script command checks that `.rhiza` reads and folds with
  no problems.
- `describe-golden.txtar` holds `describe --json`. Regenerate it when a
  request or response type changes:
  `go test ./cmd/rhiza -run TestScripts/describe-golden -update`.

## Binaries

- Build scratch binaries outside `$PATH`:
  `go build -o $TMPDIR/rhiza-x ./cmd/rhiza && $TMPDIR/rhiza-x version`.
  `~/bin/rhiza` is one slot shared by every worktree, so it may not be yours.
- Never `cp` a Go binary on arm64. Its ad-hoc signature is bound to the path
  it was linked at; a copy gets SIGKILLed (exit 137, no output). Recover one
  with `codesign --force --sign - <path>`.
- `make install` builds in place into `PREFIX` for that reason; keep it that
  way.
