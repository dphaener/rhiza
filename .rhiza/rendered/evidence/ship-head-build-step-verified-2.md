# ship-head-build-step-verified-2

evidence · `01M77BRQ4WEZNJKDXTG4F3WB6T` · actor [claude-code-a1af56f6](../actor/claude-code-a1af56f6.md)

## result

PASS

## statement

Verifier (attempt 2) at 4497a579f: `go test -trimpath ./internal/api/ -run Ship -count=1` passes every test (TestShipLandWithHeadBuild renders/lands with the named SHIP_HEAD_RHIZA binary, TestShipRenderInProcess renders in-process when none is named, the new TestShipRebaseBuildsMergedHead / TestShipRebaseRenderedConflictBuildsMergedHead show SHIP_HEAD_BUILD runs after the merge and before render, TestShipRebaseMergedHeadBuildFails shows a failed build resets the branch and pushes nothing); the package's rootguard then flags ~/code/.rhiza-root/usage.live.json, which the live daemon rewrites about once a second with real usage, unrelated to the change (a short single-test run is ok). `grep -rn isRhizaRepo` finds nothing. `rhiza workflow check ship-change` and `docs` with the commit's build: both valid (only the expected bypassed-permissions warning). The diff addresses the attempt-1 qa/quality finding (pre-merge build rendering the merged head). go vet and gofmt are clean on internal/api.

## Links

### Out

- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
