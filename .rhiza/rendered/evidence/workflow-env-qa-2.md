# workflow-env-qa-2

evidence · `01M77CPNMZXVC2Q218A6NEMD0N` · actor [claude-code-6f75672e](../actor/claude-code-6f75672e.md)

## result

PASS

## statement

QA at 3bda3170: rhiza ship rebase run by hand with a workflow-style env (SHIP_HEAD_RHIZA, SHIP_HEAD_BUILD) used both; go test -trimpath ./cmd/rhiza -run TestScripts/workflow-env$ passes and internal/api TestSessionRun (incl. the new [env] case with {run_tmp} and a plain value) passes.

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
