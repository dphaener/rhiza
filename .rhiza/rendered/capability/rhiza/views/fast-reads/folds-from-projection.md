# folds-from-projection

capability · `01M7HCGH8X1A8C8JPV1YB2X6QH` · parent [fast-reads](../fast-reads.md)

## acceptance

grep -nE 'IndexPath (!=|==) ""|newModel(Sets|Base)?\(|f\.list|FoldLog' internal/api/status.go internal/api/inbox.go internal/api/proposals.go internal/api/daemon_fold.go prints nothing; grep -nE 'IndexPath (!=|==) ""' internal/api/render.go prints nothing, and grep -c FoldLog internal/api/render.go prints 1 (the HoldsLog fallback). go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestStatus|TestInbox|TestUncoveredProposals|TestDaemonPlansProposalsPending|TestDaemonIdlePollSkipsFold|TestDaemonPollCacheInvalidates|TestDaemonFoldIndexedMatchesLog|TestRender|TestPropertyProjectionParity)' passes. Failure case: render after an entry file is edited in place still writes what the log fold renders (TestRender's in-place edit case). go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/(status|inbox-|render|daemon-plans-proposals)' passes. Touches internal/api/status.go, inbox.go (load, overlaps), proposals.go, daemon_fold.go, render.go and their tests. Must not change any --json shape or rendered file.

## statement

The ops that fold the whole log when there is no projection read only the projection: status.go's statusSets, inbox.go's inboxFold.load and overlaps (newModelSets at :690), proposals.go:107's newModelBase, daemon_fold.go:213's FoldLog and newModelBase, and render.go's empty-IndexPath branch, which keeps its strict log fold only when the index does not hold the log

## Links

### Out

- depends-on [throwaway-index-fallback](throwaway-index-fallback.md)
- decided-by [render-keeps-strict-log-fold](../../../../decision/render-keeps-strict-log-fold.md)

### In

- depends-on [interim-model-deleted](interim-model-deleted.md)
- verifies [folds-from-projection-verified](../../../../evidence/folds-from-projection-verified.md)
- verifies [folds-from-projection-qa-1](../../../../evidence/folds-from-projection-qa-1.md)
- verifies [folds-from-projection-qa-2](../../../../evidence/folds-from-projection-qa-2.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
