# log-once-is-own-leaf

decision · `01M7H101AK8J5RA1NYNEVSBXHG`

## question

Is reading the log once per call part of the projection-interface leaf or its own leaf?

## reasoning

The two are different concerns in different code (the projection in projection.go, the memo in workspace.go), each a fresh session can finish and verify alone; the orchestrator asked for one concern per leaf because reviewers flag oversized ones. The memo serves the ops that still read the log until 01M7GFTKCA deletes that path.

## statement

workspace reading the log at most once per call is its own leaf (workspace-log-once), separate from the projection interface, with no depends-on between them: workspace.log() keeps the disk read of its first call and appends the pending entries on each call, and a write to disk drops the kept read

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [reads-projection-plan-2026-11-16](../source/reads-projection-plan-2026-11-16.md)

### In

- decided-by [workspace-log-once](../capability/rhiza/views/fast-reads/workspace-log-once.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
