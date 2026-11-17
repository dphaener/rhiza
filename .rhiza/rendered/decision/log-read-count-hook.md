# log-read-count-hook

decision · `01M7HCGGT2520BYMYNB2R19HDF`

## reasoning

Every test use of readLog only counts reads, and those counts are what prove a read stays on the projection, so dropping them would lose that check. A notify-only hook keeps the check without a second way to supply the log. The hook is per Ctx, not a global counter, since global counters bumped by parallel tests made TestDaemonIdlePollSkipsFold flaky in #603.

## statement

Ctx.readLog and workspace.readLog, which replace the log read, are deleted; the tests that count log reads to show a read never decodes the log (daemon_fold_test.go:154, read_test.go:492 and :564) count through a hook on Ctx that only notes each disk read of the log and cannot change what is read

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [one-read-path-plan-2026-11-16](../source/one-read-path-plan-2026-11-16.md)

### In

- decided-by [interim-model-deleted](../capability/rhiza/views/fast-reads/interim-model-deleted.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
