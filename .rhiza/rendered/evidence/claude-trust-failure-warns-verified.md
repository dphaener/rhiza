# claude-trust-failure-warns-verified

evidence · `01M6GAE988A95GP0DD91CQ8W76` · actor [claude-code-9bbdca7d](../actor/claude-code-9bbdca7d.md)

## result

PASS

## statement

Verified at a8ca6d47 (diff 4bd51c50..HEAD: internal/api/session_run.go, session_run_trust_test.go). go test ./internal/api -run Trust -count=1: all 5 trust tests pass. TestSessionRunTrustWarnsOnUnreadableConfig writes a malformed ~/.claude.json ('{nope'), calls session.run, and asserts the warnings name the config path, the file stays byte-for-byte '{nope', and no temp file is left. The diff shows every read, parse, CreateTemp, write, chmod and rename failure returning a warning that names the file; trustFolders runs before ctx.Start in session_run.go:192 and its warning is only appended, so the agent still starts. The write goes through os.CreateTemp in the config's directory plus os.Rename, with the temp file removed on any failure. The attempt-1 finding (stale rendered file for an untracked actor) is fixed: rhiza check on a clean git archive of HEAD passes. make check: fmt, vet and every non-timing test pass. The timing budgets (slice TestBudget, api TestReadBudget) miss under load avg 10-12 on 4 cores; slice TestBudget fails the same way on base commit 4bd51c50, and internal/slice does not import internal/api, so these failures are contention, not this leaf. modelapp TestGeneratedModel and script arbiter-panes-one-window pass when re-run alone.

## Links

### Out

- verifies [claude-trust-failure-warns](../capability/rhiza/runtime/agent-adapters/claude-trust-failure-warns.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
