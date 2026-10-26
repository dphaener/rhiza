# run-status

capability · `01M5TYGK04WT4AH5KMH4RF82YD` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript reads status mid-run and after it ends.

## statement

`genesis run status [<run>]` shows a run's workflow, current step, state and the steps done.

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)

### In

- verifies [run-status-tested](../../../../evidence/run-status-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
