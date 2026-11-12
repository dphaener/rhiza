# workflow-env-qa-3

evidence · `01M77DGZJ90V1FJYSSZ5YBSSP2` · actor [claude-code-c9a49df9](../actor/claude-code-c9a49df9.md)

## result

PASS

## statement

QA at 90f39160: a scratch workflow run with [env] GOTMPDIR={run_tmp} gave its command step GOTMPDIR=<repo>/.git/rhiza/runs/<run>/tmp; TestScripts/workflow-env and internal/workflow, internal/run tests pass.

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
