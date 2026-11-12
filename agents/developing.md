# Developing Rhiza

Read this before changing Rhiza's code or running its tests. The commands
and the repo layout are in [AGENTS.md](../AGENTS.md); this doc holds the
traps that cost earlier agents retries.

## Testscripts

Scenarios live in `cmd/rhiza/testdata/script/*.txtar`.

- `$WORK` is not expanded inside single quotes. In a `stdout` regex write
  `'...'${WORK@R}'...'`.
- `stdout -count=N .` counts matching characters, not lines. Anchor a line
  regex instead.
- The builtin `chmod` takes one path per call.
- `stdin` file paths are relative to the current `cd`: use `$WORK/...`.
- Archive files get no `$WORK` expansion. A file a step reads (a workflow
  TOML) must use relative paths.
- Fake binaries (rex, claude) are archive files such as `-- bin/rex --`
  starting `#!/bin/sh`, then `chmod 755` and `env PATH=$WORK/bin:$PATH`.
  The in-process `rhiza` stays callable from them.
- When a script sets `REX_BLOCK`, the fake rex must answer `rex events`
  (exit 0, no output).
- An agent step polls every 10s by default: set `interval = "50ms"` on the
  step.
- testscript's `grep` reads a pattern starting with `-` as a flag: write
  `[-]-admin`.
- Two repos in one script can draw the same deterministic IDs: set
  `env RHIZA_TEST_SEED=<n>` around the second repo's writes.

## Scratch binaries

`$TMPDIR` is shared by every session on the VM, so another run can
overwrite `$TMPDIR/rhiza-x` mid-session. Build scratch binaries into your
session scratchpad directory, and check the binary has a flag your change
adds before trusting its output.

## RHIZA_ROOT and the run environment

In a run session `RHIZA_ROOT`, `RHIZA_CHANGE`, `RHIZA_RUN` (and
`RHIZA_ACTOR`, `RHIZA_FOR`, `RHIZA_SESSION`) point at the real repo, and
`RHIZA_ROOT` wins over the cwd. A manual `rhiza` check in a scratch repo
then hits the real one; `rhiza run start` once launched real agents this
way. Unset `RHIZA_ROOT`, `RHIZA_CHANGE` and `RHIZA_RUN` (use
`env -u ...`) for manual checks, or write a testscript with fakes. Keep
`RHIZA_ROOT` for the commands that write your own session's model entries
(`rhiza apply`, `rhiza run findings`, `rhiza session done`): your actor
record lives only there.

## GOTMPDIR

In a run, `GOTMPDIR` points under the real repo's `.git/rhiza/runs/...`, so
testscript's `$WORK` lands inside the real repo and `rhiza init` fails with
"a Rhiza repository already exists". Override it for scoped script runs:
`env -u RHIZA_ROOT GOTMPDIR=<scratchpad dir> TMPDIR=<scratchpad dir> go test ...`.

## The root daemon

A checkout under `~/code` sees the VM's real root daemon.
`internal/web` `TestOverviewReadsDaemonFromRepo` then fails ("Running, 0 of
5"), and `internal/api` can end "the test run wrote into the state folder
of the root enclosing it" because the live daemon rewrites
`.rhiza-root/usage.live.json`. Before reporting such a failure, re-run it
with `-count=1` from a scratch worktree outside `~/code`
(`git worktree add` under the scratchpad).

## rhiza render

Any landed model write (`rhiza add actor --record`, a landed change set)
leaves `.rhiza/rendered/` stale, and `make check` fails until `rhiza render`
runs and its output is committed.
