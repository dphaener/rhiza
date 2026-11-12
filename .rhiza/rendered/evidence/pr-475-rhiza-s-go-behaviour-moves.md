# pr-475-rhiza-s-go-behaviour-moves

evidence · `01M77HWP1QZMHN80X1H4SXY2VT` · actor [dphaener](../actor/dphaener.md)

## head

ad599e5f2eaf3e608038612a20dc01fe613c3ce3

## pr

https://github.com/dphaener/rhiza/pull/475

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #475 (https://github.com/dphaener/rhiza/pull/475) passed CI at ad599e5

## status

verified

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)
- verifies [rhiza-own-workflows-go-cleanup](../capability/rhiza/workflows/workflow-overrides/rhiza-own-workflows-go-cleanup.md)
- verifies [daemon-no-go-cache-trim](../capability/rhiza/orchestration/daemon/daemon-no-go-cache-trim.md)
- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)
- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
