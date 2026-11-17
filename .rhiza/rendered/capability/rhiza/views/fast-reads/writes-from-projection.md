# writes-from-projection

capability · `01M7GFZT6S9W93CH78BEGC7VCQ` · parent [fast-reads](../fast-reads.md)

## acceptance

grep -nE 'IndexPath (!=|==) ""|len\(w\.pending\) == 0|newModel\(|FoldLog' internal/api/workspace.go internal/api/read.go internal/api/change.go internal/api/ship.go internal/api/slice.go internal/api/run_model_cache.go internal/api/run_milestone.go prints nothing, and grep -n 'index.Open(ctx.Root, ctx.IndexPath)' internal/api/read.go internal/api/ship.go internal/api/repair.go prints nothing. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestWorkspace|TestDefaultView|TestRead|TestChange|TestShip|TestSlice|TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites|TestRunMilestone|TestReserveMilestone|TestWrite|TestApply|TestPropertyProjectionParity)' passes. A test shows a dry-run apply whose second write names a node the first adds resolves it through the projection's pending overlay. Failure case: a read naming a change set the log does not hold is not_found with next rhiza change ls, as today. go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/(apply|change-show|slice|ship-|repair)' passes. Touches internal/api/workspace.go, read.go, change.go, ship.go, slice.go, run_model_cache.go, run_milestone.go, repair.go and their tests. Must not delete newModel yet (interim-model-deleted does) or change any --json shape.

## statement

Every read and write of the call's own root goes through the projection, with the dry run's pending entries laid over it: workspace.model, entriesOfSet and changeSets lose their log branches, and the IndexPath and no-pending branches go from read.go's openReaderIn, change.go's runChangeShow and changeLsSets, ship.go's findShipTarget, slice.go's openSliceSource and run_model_cache.go:88, and run_milestone.go:179 reads its model from the projection instead of newModel over the disk log, and read.go:183, ship.go:398 and repair.go:161 take the open index from the projection (its store()) instead of calling index.Open on ctx.IndexPath

## Links

### Out

- depends-on [throwaway-index-fallback](throwaway-index-fallback.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)

### In

- depends-on [interim-model-deleted](interim-model-deleted.md)
- verifies [writes-from-projection-verified](../../../../evidence/writes-from-projection-verified.md)
- verifies [writes-from-projection-qa-1](../../../../evidence/writes-from-projection-qa-1.md)
- verifies [writes-from-projection-qa-2](../../../../evidence/writes-from-projection-qa-2.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
