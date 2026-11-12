# ship-head-build-step-verified

evidence · `01M779B38GMYFZY4HNNDHZNP17` · actor [claude-code-0473f66f](../actor/claude-code-0473f66f.md)

## result

PASS

## statement

Verified at 9bdbe8da6. go test -trimpath -count=1 ./internal/api -run Ship passes: TestShipLandWithHeadBuild (SHIP_HEAD_RHIZA names a build; land and render run with that binary, 2 exec calls, no go build, in-process landed nothing), TestShipLandHeadBuildMissing (a named but missing build fails naming SHIP_HEAD_RHIZA, nothing run or pushed), TestShipLandHeadBuildErrorDetail, and TestShipLand (no head build named: lands and renders in-process). isRhizaRepo, goModule and rhizaModule are gone from the tree (grep finds none); shipHeadBuild reads only SHIP_HEAD_RHIZA. .rhiza/workflows/ship-change.toml is the built-in plus [env] SHIP_HEAD_RHIZA={run_tmp}/rhiza-head and head-build-rebase/land/rebase-merge steps (go build -o $SHIP_HEAD_RHIZA ./cmd/rhiza) that are the only routes into rebase, land and rebase-merge, the only steps that call shipRender/shipHeadBuild; docs.toml mirrors it and bugfix.toml sets the env. rhiza workflow check (a build of this commit) reports valid, exit 0, for bugfix, docs, docs-build, docs-review, plan-bug and ship-change; the only warnings are the expected bypassed-permissions ones.

## Links

### Out

- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
