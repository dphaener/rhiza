# projection-interface

capability · `01M7GFYXY7FTRWHX1NNNSW8BWA` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestProjection|TestPropertyProjectionParity|TestWrite)' passes, with TestPropertyProjectionParity comparing the projection's views (not writeIndex.model) against a fresh strict fold. New tests in internal/api/projection_test.go: a dry-run write with pending entries sees them through the projection (its change sets, entries and model) and leaves the index file unchanged (its Stamp before and after equal); one call that asks the projection for the model, the change sets and a change set's entries three times opens the index once (counted through a test hook on the opener); with ctx.IndexPath empty the workspace's projection is nil and the op answers from the log as before. Failure case: an index file that is not a database is rebuilt and the call answers as before; an index that cannot be opened returns that error, not a silent fold. Touches internal/api/projection.go (new), write_indexed.go, workspace.go, run_model_cache.go, inbox.go only for inboxGateRead's opener, and their tests. Must not move status, the inbox's fold, render, change ls or the daemon fold onto it (their own leaves do), must not change read.go's reader, and must not delete the interim model or an IndexPath branch (01M7GFTKCA does).

## statement

internal/api has one projection interface in a new internal/api/projection.go (change sets with their states and folds, a change set's entries, the stored problems, the model of the current change set's view or of as-is, and the dry run's pending entries laid over all of them), whose only implementation is writeIndex's index-plus-pending code moved behind it; workspace opens it lazily, at most once per call, when ctx.IndexPath is set, and the writes, run_model_cache.go and inboxGateRead that build a writeIndex today get it through the interface

## Links

### Out

- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)
- decided-by [projection-grows-from-write-index](../../../../decision/projection-grows-from-write-index.md)

### In

- depends-on [status-from-projection](status-from-projection.md)
- depends-on [inbox-from-projection](inbox-from-projection.md)
- depends-on [render-from-projection](render-from-projection.md)
- depends-on [change-reads-from-projection](change-reads-from-projection.md)
- depends-on [daemon-fold-from-projection](daemon-fold-from-projection.md)
- depends-on [slice-and-show-from-projection](slice-and-show-from-projection.md)
- verifies [projection-interface-verified](../../../../evidence/projection-interface-verified.md)
- verifies [projection-interface-qa-1](../../../../evidence/projection-interface-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
