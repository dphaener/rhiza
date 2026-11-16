# index-as-is-model-verified

evidence · `01M7H4JHHKM56EA7M8EE4ZZX3S` · actor [claude-code-8a23b3f6](../actor/claude-code-8a23b3f6.md)

## result

PASS

## statement

Verified commit ca9251a3f. Re-ran its Checks: go test -trimpath -count=1 -p 2 ./internal/index -run '^(TestAsIs|TestPropertyIncremental)' exit 0 (TestAsIsProperty: 1055 steps, 53 logs with two landed heads, 46 with a duplicate handle, 82 with a rebuild; TestAsIsEmpty, TestAsIsClosed, TestPropertyIncremental pass); go test ... ./internal/api -run TestPropertyProjectionParity exit 0; go vet ./internal/index exit 0; rhiza ship generate exit 0 (0 files); rhiza render exit 0. Diff read: asis.go's (*Index).AsIs loads the stored as-is rows via loadSeed (nodes with fields/history/aspects, links, as-is problems from the problems table) and finishes them with model.ToBeFrom(seed, nil); it decodes no log entry and runs no fold. The property test draws logs with indextest (plus an injected landed duplicate handle), catches up step by step and compares model nodes, links, tree and as-is problems with model.AsIs over changeset.FoldLog after every step, failing a vacuous run. Empty log gives empty model, nil problems, no error; closed index errors. Only internal/index/asis.go, asis_test.go and .rhiza files changed: no stored row, schema or internal/model rule touched; no TODOs.

## Links

### Out

- verifies [index-as-is-model](../capability/rhiza/views/fast-reads/index-as-is-model.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
