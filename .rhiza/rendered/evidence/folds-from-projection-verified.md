# folds-from-projection-verified

evidence · `01M7K6EHJNAJKQJQFRV5YDQGJW` · actor [claude-code-c481df8f](../actor/claude-code-c481df8f.md)

## result

PASS

## statement

Detached checkout at 22f10f8db. The acceptance grep over status.go, inbox.go, proposals.go and daemon_fold.go printed nothing; the IndexPath grep over render.go printed nothing and grep -c FoldLog render.go printed 1 (the !HoldsLog fallback). go vet ./internal/api exit 0. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestStatus|TestInbox|TestUncoveredProposals|TestDaemonPlansProposalsPending|TestDaemonIdlePollSkipsFold|TestDaemonPollCacheInvalidates|TestDaemonFoldIndexedMatchesLog|TestRender|TestPropertyProjectionParity)' ok, including TestRenderIndexedAfterInPlaceEdit (the in-place edit failure case) and the parity test. go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/(status|inbox-|render|daemon-plans-proposals)' ok with RHIZA_/SHIP_/REX_ vars unset. Diff read: only the five named files change in code; the log-fold branches are removed rather than left as fallbacks, the projection is never nil (useProjection opens a throwaway index when IndexPath is empty, a failedProjection on error), no --json shape or rendered file changes, no stubs or TODOs.

## Links

### Out

- verifies [folds-from-projection](../capability/rhiza/views/fast-reads/folds-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
