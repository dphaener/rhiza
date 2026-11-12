# ship-head-build-step-qa-3

evidence · `01M77DGZE01018N5BYXEDK1DYS` · actor [claude-code-c9a49df9](../actor/claude-code-c9a49df9.md)

## result

PASS

## statement

QA at 90f39160 with a -trimpath build: in a scratch Rhiza root, an outer workflow including an inner one whose head-build step runs sh -c "${SHIP_HEAD_BUILD:?...}" (bugfix's shape) succeeded using the outer [env]'s SHIP_HEAD_BUILD (the inner's was ignored) with GOTMPDIR the run's tmp; with no SHIP_HEAD_BUILD in the outer [env] the step failed with 'SHIP_HEAD_BUILD: is not set' and the run failed. bugfix.toml's SHIP_HEAD_BUILD built a working rhiza; rhiza workflow check passes for ship-change, bugfix, docs, plan-bug, docs-build and docs-review. TestScripts/ship-land-branch-build (now SHIP_HEAD_RHIZA-based) and internal/api Ship|SessionRun|Daemon tests pass (run outside ~/code).

## Links

### Out

- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
