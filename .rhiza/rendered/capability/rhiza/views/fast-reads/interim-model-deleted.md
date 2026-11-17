# interim-model-deleted

capability · `01M7GFZTBT1GB2DR05VA9EGCSV` · parent [fast-reads](../fast-reads.md)

## acceptance

grep -nE 'func newModel|newModelSets|newModelBase|modelBase|strictVersions|readLog' internal/api/*.go | grep -v _test.go prints nothing, and grep -nE '\.readLog\b|newModel(Sets|Base)?\(' internal/api/*_test.go prints nothing (the read_test.go fixture func readLog(t) and its readLog(t) calls keep their name). grep -nE 'IndexPath (!=|==) ""' internal/api/*.go | grep -v _test.go prints at most the one line of the projection open's throwaway-index check (projection.go or write_indexed.go, from throwaway-index-fallback); tests may still compare IndexPath with "" (slice_test.go:291, change_indexed_test.go:225). go vet ./internal/api is clean; go test -trimpath -count=1 -p 2 ./internal/api passes; go test -trimpath -count=1 -p 2 ./internal/api -run '^TestPropertyProjectionParity$' -rapid.checks=1000 passes. TestDaemonIdlePollSkipsFold and the read_test.go read counts still show their reads decode no log, through the hook. go test -list '.*' ./... | grep -c '^Test' and ls cmd/rhiza/testdata/script/*.txtar | wc -l are no lower than at a98fd4051 (1,948 and 398); a deleted test is named in the commit message with the parity case or test that now covers it. Failure case: a call with a damaged index file still rebuilds it and answers (TestProjectionIndexNotADatabase), with no fold fallback left to hide a projection bug. Touches internal/api/model.go, api.go, workspace.go, parity_test.go, the tests that set Ctx.readLog or workspace.readLog (daemon_fold_test.go, read_test.go, workspace_test.go) and the _test.go files of the deleted code. Must not change any --json shape or rendered file (make self-check passes).

## statement

The interim model is gone: internal/api/model.go loses newModel, newModelSets, newModelBase, the modelBase type with its model and conflicts methods, strictVersions, head, copyNode, copyLink, ambiguousNodes, and the model struct with its methods (sortedNodes, sortedLinks, Resolve, Node, ResolveLink, LinkBetween, LinksTo, LinksFrom, DependsOn, Children, HandleTaken, ActorByEmail), while the Model interface, Node, Link, dependencyOrder (write_indexed.go:240) and ambiguousLinks (write_indexed.go:441) stay; Ctx.readLog and workspace.readLog go (replaced by a notify-only log-read hook for the tests that count reads); the tests of the deleted code only are deleted; and the parity test compares the projection with a fresh strict fold alone

## Links

### Out

- depends-on [writes-from-projection](writes-from-projection.md)
- depends-on [other-roots-own-projection](other-roots-own-projection.md)
- depends-on [throwaway-index-fallback](throwaway-index-fallback.md)
- depends-on [folds-from-projection](folds-from-projection.md)
- decided-by [projection-targets](../../../../decision/projection-targets.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)
- decided-by [log-read-count-hook](../../../../decision/log-read-count-hook.md)

### In

- verifies [interim-model-deleted-verified](../../../../evidence/interim-model-deleted-verified.md)
- verifies [interim-model-deleted-qa-1](../../../../evidence/interim-model-deleted-qa-1.md)
- verifies [interim-model-deleted-qa-2](../../../../evidence/interim-model-deleted-qa-2.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
