# workflow-env

capability · `01M7642MCFVAH4VXGAXZNTW4S7` · parent [workflow-file](../workflow-file.md)

## acceptance

A testscript runs a workflow setting FOO={run_tmp} and shows a command step sees it and the folder is gone after the run; a workflow without env sees no GOTMPDIR set by Rhiza.

## statement

A workflow file may set environment variables for every process its steps start, with {run_tmp} filled in as the run's own temp folder, which Rhiza creates when a step starts and removes when the run ends; Rhiza itself no longer sets GOTMPDIR (run-private-gotmpdir is replaced by this and retired).

## Links

### In

- verifies [workflow-env-verified](../../../../../evidence/workflow-env-verified.md)
- verifies [workflow-env-qa-1](../../../../../evidence/workflow-env-qa-1.md)
- verifies [workflow-env-verified-2](../../../../../evidence/workflow-env-verified-2.md)
- verifies [workflow-env-qa-2](../../../../../evidence/workflow-env-qa-2.md)
- verifies [workflow-env-qa-3](../../../../../evidence/workflow-env-qa-3.md)
- verifies [pr-475-rhiza-s-go-behaviour-moves](../../../../../evidence/pr-475-rhiza-s-go-behaviour-moves.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
