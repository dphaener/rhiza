# run-resume

capability · `01M5TYGK7EFHKPB9JG3A2MH2G2` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript stops a run after step one, resumes it, and step one runs once.

## statement

`genesis run resume <run>` continues a stopped run from its journal without a daemon, never redoing a finished step; it refuses an interrupted run whose change set has landed, been rejected, withdrawn or superseded.

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)

### In

- depends-on [manual-resume-refuses-unaccepted](manual-resume-refuses-unaccepted.md)
- verifies [run-resume-tested](../../../../evidence/run-resume-tested.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
