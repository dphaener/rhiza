# skipped-run-in-inbox

decision · `01M62RH250PMR7VXXN022JN5N0`

## statement

An interrupted run whose change set is no longer accepted is not resumed and waits in the inbox with the change set's state and the next step, rhiza run abandon

## wording

Show in inbox, the agent's option, picked by the owner

## Links

### Out

- derived-from [resume-only-in-flight-plan-interview](../source/resume-only-in-flight-plan-interview.md)

### In

- decided-by [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)
- decided-by [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
