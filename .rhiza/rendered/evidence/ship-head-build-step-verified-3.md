# ship-head-build-step-verified-3

evidence · `01M77D3NC9P1MAF6Y60214GDGK` · actor [claude-code-b7a6cc99](../actor/claude-code-b7a6cc99.md)

## result

PASS

## statement

Verifier (attempt 5) at a4395c22b: read the diff. bugfix.toml now sets SHIP_HEAD_BUILD. The head-build-* steps in ship-change.toml and docs.toml run sh -c "${SHIP_HEAD_BUILD:?...}"; a shell check shows this exits 2 when the variable is unset and runs the command when it is set. ship-land-branch-build.txtar was rewritten for SHIP_HEAD_RHIZA, and no isRhizaRepo is left in any Go file. Ran go test -trimpath ./internal/api -run Ship: ok. TestShipLandWithHeadBuild, TestShipLandHeadBuildMissing, TestShipRenderInProcess, TestShipRebaseBuildsMergedHead, TestShipRebaseRenderedConflictBuildsMergedHead, TestShipRebaseMergedHeadBuildFails and TestShipLandHeadBuildErrorDetail all PASS. go test -trimpath ./cmd/rhiza -run 'TestScripts/ship-land-branch-build$': ok. go test -trimpath ./internal/workflow/...: ok. With a scratch build, rhiza workflow check ship-change, docs and bugfix each report valid (only the expected bypassed-permissions warning). make self-check fails only on render files for untracked run-log actors in this checkout, not on anything the commit touches.

## Links

### Out

- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
