# run-start

capability · `01M5TYGJASX1Q1J7DV6DVHB0DZ` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript starts a run and finds it in the journal.

## statement

`genesis run start <change set>` starts a run of a workflow on that change set and prints its run ID.

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)
- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)

### In

- verifies [run-start-tested](../../../../evidence/run-start-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
