# run-resume

capability · `01M5TYGK7EFHKPB9JG3A2MH2G2` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript stops a run after step one, resumes it, and step one runs once.

## statement

`genesis run resume <run>` continues a stopped run from its journal without a daemon, never redoing a finished step.

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)

### In

- verifies [run-resume-tested](../../../../evidence/run-resume-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
