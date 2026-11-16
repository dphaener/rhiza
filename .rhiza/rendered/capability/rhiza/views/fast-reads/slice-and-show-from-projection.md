# slice-and-show-from-projection

capability · `01M7H1E0YWEKQ7YB02B4KH9WSH` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestSlice|TestChangeShow)' passes; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(slice|change-show)' passes for every slice*.txtar and change-show*.txtar script. grep -n 'index.Open' internal/api/change_indexed.go prints nothing, and in internal/api/slice.go openIndex is called only on the empty-IndexPath branch. A test with ctx.IndexPath set counts index opens through projection-interface's opener hook and shows slice and change show each open the index once, through the projection, with the same answers as before. Failure case: slice outside a repository is still not_found with next rhiza init, and change show with no change set named or current is still the same usage error. Touches internal/api/slice.go, change_indexed.go and their tests. Must not change slice's or change show's text or --json output, nor the answers indexModel gives, and must not delete openIndex (01M7GFTKCA does).

## statement

rhiza slice (internal/api/slice.go:107-140, its own openIndex and index.Open) and rhiza change show over the index (runChangeShowIndexed, internal/api/change_indexed.go:217, its own index.Open) take their index reads (change sets, a change set's entries, the view's model) from the workspace's projection when ctx.IndexPath is set instead of opening internal/index themselves; with ctx.IndexPath empty, slice keeps its openIndex of the default path, which 01M7GFTKCA's throwaway-index-fallback removes

## Links

### Out

- depends-on [projection-interface](projection-interface.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- verifies [slice-and-show-from-projection-verified](../../../../evidence/slice-and-show-from-projection-verified.md)
- verifies [slice-and-show-from-projection-qa-1](../../../../evidence/slice-and-show-from-projection-qa-1.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
