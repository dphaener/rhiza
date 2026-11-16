# daemon-fold-from-projection

capability · `01M7GFYY98FH6TH6DY2EZTKFPR` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'Daemon' passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/daemon' passes for every daemon*.txtar script. A test with ctx.IndexPath set and Ctx.readLog counting, over a log with a proposed change set with leaves, one without, an accepted one and overlapping in-flight change sets, shows the poll reads the log zero times and its candidates, proposed list, plan skips, overlaps and lifecycle moves equal those over the log. Failure case: a stamp change from a vanished entry file (the index rebuilds) still gives the same candidates on the next poll. Touches internal/api/daemon_fold.go and its tests. Must not change the daemon's schedule, its journal records or any --json shape, and must not delete newModelBase (01M7GFTKCA does).

## statement

The daemon's poll fold (internal/api/daemon_fold.go:197) takes its change sets from the projection instead of changeset.FoldLog over the whole log, and the models it reads per candidate, for the proposed list and for the overlaps and actors (:215-230) from the projection's views instead of newModelBase over the folded list

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- verifies [daemon-fold-from-projection-verified](../../../../evidence/daemon-fold-from-projection-verified.md)
- verifies [daemon-fold-from-projection-qa-1](../../../../evidence/daemon-fold-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
