# status-from-projection

capability · `01M7GFYY0KNAMKMJK1MHE68E99` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run Status passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/status' passes for every status*.txtar script. A test with ctx.IndexPath set and Ctx.readLog counting shows status, and status with Hook, read the log zero times and give the same counts and change set as over the log. Measure in the commit's Checks: block, each run wrapped as /home/rhiza/orchestrator/machine-rules.md says: go build -trimpath -o $TMPDIR/rhiza-x ./cmd/rhiza, then in this repo's main checkout run $TMPDIR/rhiza-x status once and /usr/bin/time -f '%e s' $TMPDIR/rhiza-x status, within 200 ms (2026-11-16: 1.93 s). Failure case: status --hook outside a repo stays silent with exit 0, and with an index that cannot be opened it still exits 0. Touches internal/api/status.go and its tests. Must not change status's text or --json output.

## statement

rhiza status and its session-start hook (StatusRequest.Hook) count change sets by state from the projection instead of folding the whole log (internal/api/status.go:123), so the hook paid on every session start costs a few index rows

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- decided-by [projection-targets](../../../../decision/projection-targets.md)

### In

- verifies [status-from-projection-verified](../../../../evidence/status-from-projection-verified.md)
- verifies [status-from-projection-qa-1](../../../../evidence/status-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
