# daemon-repo-session-verified

evidence · `01M6FKR9ZYZDYTBD3X7MMHXKW2` · actor [claude-code-036f0c82](../actor/claude-code-036f0c82.md)

## result

PASS

## statement

Verifier at 6f0d4416: go test ./cmd/rhiza -run TestScripts/daemon -count=1 passes all 12 daemon scripts, including daemon-repo-session.txtar, which runs rhiza daemon --once against a fake rex with REX_SESSION unset: with rex session ls empty it logs exactly one 'rex session new rhiza-repo' and the run start tab as 'rex window new -s rhiza-repo', with no window new lacking -s rhiza-repo; with rhiza-repo listed it logs no session new and every window new carries -s rhiza-repo; with REX_SESSION=session:other the tabs still carry -s rhiza-repo and session:other never appears. go test ./internal/api passes, and make check exits 0. Diff read: ensureRepoSession (internal/api/daemon.go) is called before the run, resume and plan-request window new calls, names the session rhiza-<base of the first git worktree>, and overrides ctx.RexSession. A read-only 'rex session ls --json' on the real server shows the sessions[].label shape the code parses, and 'rex session new --help' and 'rex window new --help' confirm the LABEL, -c and -s arguments it passes.

## Links

### Out

- verifies daemon-repo-session (retired)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
