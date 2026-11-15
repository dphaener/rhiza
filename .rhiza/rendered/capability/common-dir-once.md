# common-dir-once

capability · `01M7FT14KGB9BV6KWFB238MDPT`

## acceptance

A new test with a counting Git (go test -trimpath -p 2 -count=1 ./internal/api -run '^TestGitCommonDirOnce$') shows two gitCommonDir calls on one Ctx making one git call, rev-parse --path-format=absolute --git-common-dir, and returning the same absolute, cleaned path, from the main checkout and from a worktree, and orphanLogEntries after a gitCommonDir call on the same Ctx making no further --git-common-dir call; a second test (go test -trimpath -p 2 -count=1 ./internal/api -run '^TestOrphanLogSymlinkRoot$') with a real git repo whose main checkout is reached through a symlinked directory (Ctx.Root the symlinked path) shows the kept common dir equal to rev-parse --absolute-git-dir and orphanLogEntries still listing an untracked orphan log entry; go test -trimpath -p 2 ./internal/api -run 'Session|Orphan|Untracked' and go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/^(session-run|session-start|apply|run-amend|check-orphan-log|model-commit-untracked)$' pass, IDs unchanged (the scripts' deterministic IDs still match). Measure: strace -f -e trace=execve on go test -run 'TestScripts/^run-ship-conflict-fix-more$' ./cmd/rhiza counts at most one rev-parse --git-common-dir per rhiza child, rhiza check children included (2026-11-15: 757 for 287). Touches internal/api/session.go, internal/api/api.go (Ctx), internal/api/check.go (orphanLogEntries), cmd/rhiza/ctx.go and their tests.

## statement

A rhiza process runs git rev-parse --git-common-dir at most once for its checkout, with --path-format=absolute so git resolves symlinks: the test-mode ID offset, every gitCommonDir call and rhiza check's orphan-log test (orphanLogEntries) reuse the first answer, and orphanLogEntries still finds orphans in a main checkout reached through a symlink.

## Links

### Out

- decided-by [common-dir-asked-once](../decision/common-dir-asked-once.md)

### In

- verifies [common-dir-once-verified](../evidence/common-dir-once-verified.md)
- verifies [common-dir-once-qa-1](../evidence/common-dir-once-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
