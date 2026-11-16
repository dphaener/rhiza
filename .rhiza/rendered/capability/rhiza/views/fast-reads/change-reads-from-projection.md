# change-reads-from-projection

capability · `01M7GFYY75HKKZF3E9JKJBT985` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'ChangeLs' passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/change' passes for every change*.txtar script. A test with ctx.IndexPath set and Ctx.readLog counting shows change ls reads the log zero times and lists the same change sets and problems as over the log, a damaged entry's problem included. Failure case: an entry file that vanished after the index saw it (a rebase) makes the projection rebuild, and change ls lists the remaining change sets with no stale one. Touches internal/api/change.go and its tests. Must not change change ls's text or --json output, nor change show, whose index path change_indexed.go already has.

## statement

rhiza change ls (internal/api/change.go:669) takes the change sets, their kind, state, intent and title, and the log and fold problems it lists from the projection (its folded change sets and stored problems) instead of workspace.log() and changeset.FoldLog over the whole log

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- verifies [change-reads-from-projection-verified](../../../../evidence/change-reads-from-projection-verified.md)
- verifies [change-reads-from-projection-qa-1](../../../../evidence/change-reads-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
