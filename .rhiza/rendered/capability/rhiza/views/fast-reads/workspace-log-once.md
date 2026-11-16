# workspace-log-once

capability · `01M7H101RT64ATZVGM94276F58` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestWorkspaceLog|TestChange|TestWrite)' passes. New tests in internal/api/workspace_test.go, counting disk reads through Ctx.readLog: an op over the log path (ctx.IndexPath empty) that calls model, changeSets and entriesOfSet several times reads the disk once; a dry run of several writes still chains, each later read seeing the earlier pending entries, with one disk read; a real write followed by a read in the same workspace reads the disk again and sees the written entry. Failure case: a log with a damaged entry gives the same problems on the second call as on the first. Touches internal/api/workspace.go and its tests only. Must not change any op's output.

## statement

workspace.log() (internal/api/workspace.go:146) reads the log from disk at most once per workspace, keeping that read and laying the dry run's pending entries on it at each call, and a write to disk drops the kept read so the next call sees the new entry

## Links

### Out

- decided-by [log-once-is-own-leaf](../../../../decision/log-once-is-own-leaf.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- verifies [workspace-log-once-verified](../../../../evidence/workspace-log-once-verified.md)
- verifies [workspace-log-once-qa-1](../../../../evidence/workspace-log-once-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
