# run-journal-file

capability · `01M5TYGJNH9AMN4S8C5QSY3CW4` · parent [workflow-runs](../workflow-runs.md)

## acceptance

After a run the journal holds every step and `.genesis/log/` is unchanged.

## statement

Each run keeps a local journal of its steps, outcomes and visits outside the model's log.

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)
- decided-by [work-state-outside-log](../../../../decision/work-state-outside-log.md)

### In

- verifies [run-journal-file-tested](../../../../evidence/run-journal-file-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
