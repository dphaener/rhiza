# index-as-is-model

capability · `01M7H101ZYJWDMD0XH748H3H6C` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/index -run '^(TestAsIs|TestPropertyIncremental)' passes. A new property test in internal/index, drawing logs with internal/index/indextest (two landed heads, rebuilds, a duplicate handle), catches an index up step by step and checks that AsIs's model and problems equal model.AsIs over changeset.FoldLog of the same log after every step. Failure case: an empty log gives an empty model and no error; a closed index returns an error. Touches a new internal/index/asis.go and its test only. Must not change a stored row, the schema or any rule in internal/model.

## statement

internal/index has a read, (*Index).AsIs() (model.Model, []model.Problem, error), that assembles the stored as-is rows (every node with its fields and history, every link, the as-is problems) into an internal/model Model without decoding a log entry

## Links

### Out

- decided-by [index-as-is-read](../../../../decision/index-as-is-read.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)

### In

- depends-on [render-from-projection](render-from-projection.md)
- verifies [index-as-is-model-verified](../../../../evidence/index-as-is-model-verified.md)
- verifies [index-as-is-model-qa-1](../../../../evidence/index-as-is-model-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
