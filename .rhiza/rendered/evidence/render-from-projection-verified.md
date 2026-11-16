# render-from-projection-verified

evidence · `01M7H546GR91YW6B41EK6Q06R0` · actor [claude-code-e9dd4da5](../actor/claude-code-e9dd4da5.md)

## result

PASS

## statement

Verified da24c315c under the capped scope: go test -trimpath -count=1 -p 2 ./internal/api -run Render passed (TestRenderIndexedMatchesLog renders over the index with an Evidence node's creator and a retired link endpoint, then over the log, every file unchanged); TestPropertyProjectionParity passed; go test ./cmd/rhiza -run TestScripts/render passed render, render-generate and render-refuses-damaged-log; make self-check in a detached checkout of the commit exit 0 with git status --short .rhiza/rendered empty; the parent commit's build rendered that checkout with 0 written (byte-identical). Diff: render.go calls damagedLog unchanged before renderInput, renderInput keeps the empty-IndexPath FoldLog branch, projection gains asIs; internal/api/model.go untouched; no TODOs or stubs.

## Links

### Out

- verifies [render-from-projection](../capability/rhiza/views/fast-reads/render-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
