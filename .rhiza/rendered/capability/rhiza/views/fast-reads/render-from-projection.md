# render-from-projection

capability · `01M7GFYY4Z3CRG97YV48H46RQN` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run Render passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/render' passes for every render*.txtar script; make self-check passes, so rhiza render on this repo writes byte-identical .rhiza/rendered/ files (git status --short .rhiza/rendered prints nothing after it). A test with ctx.IndexPath set renders the same files as over the log, including an Evidence node's creator and a retired link endpoint. Failure case: a log with an entry that does not decode still makes render refuse with the same error as today (render-refuses-damaged-log). Measure in the commit's Checks: block, wrapped as /home/rhiza/orchestrator/machine-rules.md says: /usr/bin/time -f '%e s' $TMPDIR/rhiza-x render in this repo's main checkout, before and after. Touches internal/api/render.go and its tests. Must not change a rendered file's content or damagedLog's rules.

## statement

rhiza render takes the change sets and the as-is model it renders from the projection (its folded change sets and index-as-is-model's AsIs) instead of changeset.FoldLog and model.AsIs over the whole log (internal/api/render.go:99-100), and keeps damagedLog's own read so it still refuses a damaged log

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- depends-on [index-as-is-model](index-as-is-model.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [render-keeps-damage-read](../../../../decision/render-keeps-damage-read.md)

### In

- verifies [render-from-projection-verified](../../../../evidence/render-from-projection-verified.md)
- verifies [render-from-projection-qa-1](../../../../evidence/render-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
