# plan-run-resumes-landed-record

capability · `01M6GD3GB66TDYWWHF23GMTFZE` · parent [plan-gate](../plan-gate.md)

## acceptance

A testscript resumes an interrupted plan run whose record change set has landed and the run goes on past its stopped step; in the same suite rhiza run resume still exits non-zero naming the state and rhiza run abandon for a plan run whose change set is withdrawn and for a build run whose change set has landed (run-resume-unaccepted.txtar keeps passing), each journal unchanged

## statement

rhiza run resume continues an unfinished plan-workflow run whose change set has landed (a plan run lands its record during the interview), while it still refuses a plan run whose change set is rejected, withdrawn or superseded, and any other workflow's run whose change set is no longer in flight

## Links

### Out

- depends-on [plan-run-stops-open](plan-run-stops-open.md)
- decided-by [plan-resume-landed-record-milestone-6](../../../../../../decision/plan-resume-landed-record-milestone-6.md)

### In

- verifies [plan-run-resumes-landed-record-verified](../../../../../../evidence/plan-run-resumes-landed-record-verified.md)
- verifies [plan-run-resumes-landed-record-qa-1](../../../../../../evidence/plan-run-resumes-landed-record-qa-1.md)
- verifies [plan-run-resumes-landed-record-verified-2](../../../../../../evidence/plan-run-resumes-landed-record-verified-2.md)
- verifies [pr-216-fix](../../../../../../evidence/pr-216-fix.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
