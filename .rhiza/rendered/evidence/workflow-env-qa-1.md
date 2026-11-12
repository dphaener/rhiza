# workflow-env-qa-1

evidence · `01M77AP4TS8SNXCBY96DQHYB0P` · actor [claude-code-174ed9c7](../actor/claude-code-174ed9c7.md)

## result

PASS

## statement

QA at 46f1cad8 with a -trimpath build, in a scratch root and repo with RHIZA_* and GOTMPDIR unset: a workflow with [env] GOTMPDIR={run_tmp}, NESTED={run_tmp}/sub, PLAIN=x ran with rhiza run start --foreground and succeeded; its command step saw GOTMPDIR=<repo>/.git/rhiza/runs/<run>/tmp (a directory while the step ran), NESTED=<that>/sub and PLAIN=x; its [[cleanup]] step saw the same GOTMPDIR with the folder still present; after the run the folder was gone. A workflow without [env] saw no GOTMPDIR. A workflow setting RHIZA_CHANGE in [env] fails rhiza workflow check with 'env: RHIZA_CHANGE is Rhiza's own variable and cannot be set'. go test ./internal/workflow/... ./internal/run/... and TestScripts/workflow-env pass.

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
