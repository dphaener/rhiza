# inbox-from-projection

capability · `01M7GFYY2VQJHCVC83XPQFR2EA` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run Inbox passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/inbox' passes for every inbox*.txtar script. A test with ctx.IndexPath set and Ctx.readLog counting shows an inbox with overlapping in-flight change sets and a run waiting on questions reads the log zero times and lists the same items as over the log. grep -n 'index.Open\|openWriteIndex' internal/api/inbox.go prints nothing. Failure case: a run whose worktree is gone or whose log does not read gives the inbox item without its questions, as today, and no error. Touches internal/api/inbox.go and its tests. Must not change the inbox's items or --json shape, and must not change internal/run.

## statement

rhiza inbox's inboxFold (internal/api/inbox.go:602-619) takes change set states, the in-flight change sets and the overlaps' as-is model (:685) from the workspace's projection instead of opening internal/index itself, folding the log only when the projection is nil, and inboxGateRead (:1405) reads the run worktree through that worktree's projection instead of building a writeIndex by hand

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- decided-by [projection-targets](../../../../decision/projection-targets.md)
- decided-by [inbox-rescoped-after-parity](../../../../decision/inbox-rescoped-after-parity.md)

### In

- verifies [inbox-from-projection-verified](../../../../evidence/inbox-from-projection-verified.md)
- verifies [inbox-from-projection-qa-1](../../../../evidence/inbox-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
