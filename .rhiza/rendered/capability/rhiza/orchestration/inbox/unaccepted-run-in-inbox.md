# unaccepted-run-in-inbox

capability · `01M62RH3PBJ12H0CYSMCMV7KAS` · parent [inbox](../inbox.md)

## acceptance

A test with an interrupted run on a landed change set lists it in rhiza inbox with the state and the abandon command, and the daemon does not resume it

## statement

An interrupted run on a change set no longer accepted waits in the inbox, naming the change set's state and rhiza run abandon as the next step

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)

### In

- verifies [unaccepted-run-in-inbox-verified](../../../../evidence/unaccepted-run-in-inbox-verified.md)
- verifies [unaccepted-run-in-inbox-qa-1](../../../../evidence/unaccepted-run-in-inbox-qa-1.md)
- verifies [unaccepted-run-in-inbox-verified-2](../../../../evidence/unaccepted-run-in-inbox-verified-2.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
