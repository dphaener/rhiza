# daemon-resumes-accepted-only

capability · `01M62RH3705M0GA17WVVXKBGM8` · parent [daemon](../daemon.md)

## acceptance

A test resumes an interrupted run on an accepted change set, and leaves interrupted runs on landed, rejected, withdrawn, proposed and draft change sets unresumed

## statement

The daemon resumes an interrupted run only while its change set is accepted

## Links

### Out

- depends-on [daemon-resume-once](daemon-resume-once.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)

### In

- verifies [daemon-resumes-accepted-only-verified](../../../../evidence/daemon-resumes-accepted-only-verified.md)
- verifies [daemon-resumes-accepted-only-qa-1](../../../../evidence/daemon-resumes-accepted-only-qa-1.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
