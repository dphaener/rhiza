# manual-resume-refuses

decision · `01M62RH2N6XDVSJV0K31Z999SJ`

## statement

rhiza run resume refuses a run whose change set is no longer accepted, as the daemon does, except a plan-workflow run whose change set has landed (Decision plan-resume-landed-record-milestone-6)

## wording

Refuse too, the agent's option, picked by the owner

## Links

### Out

- derived-from [resume-only-in-flight-plan-interview](../source/resume-only-in-flight-plan-interview.md)

### In

- decided-by [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
