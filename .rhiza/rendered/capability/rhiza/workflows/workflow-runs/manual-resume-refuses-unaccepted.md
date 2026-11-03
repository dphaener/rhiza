# manual-resume-refuses-unaccepted

capability · `01M62RH4605WT6153H4AZSRG87` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A test runs rhiza run resume on an interrupted run whose change set is withdrawn, and it exits non-zero naming the state and rhiza run abandon, leaving the journal unchanged

## statement

rhiza run resume refuses an interrupted run whose change set is no longer accepted, naming the state and rhiza run abandon; a plan-workflow run whose change set has landed is the exception (plan-run-resumes-landed-record)

## Links

### Out

- depends-on [run-resume](run-resume.md)
- decided-by [manual-resume-refuses](../../../../decision/manual-resume-refuses.md)

### In

- verifies [manual-resume-refuses-unaccepted-verified](../../../../evidence/manual-resume-refuses-unaccepted-verified.md)
- verifies [manual-resume-refuses-unaccepted-qa-1](../../../../evidence/manual-resume-refuses-unaccepted-qa-1.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-216-fix](../../../../evidence/pr-216-fix.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
