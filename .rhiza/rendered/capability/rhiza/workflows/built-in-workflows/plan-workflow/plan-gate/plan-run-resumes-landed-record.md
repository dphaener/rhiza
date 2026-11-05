# plan-run-resumes-landed-record

capability · `01M6GD3GB66TDYWWHF23GMTFZE` · parent [plan-gate](../plan-gate.md)

## acceptance

A testscript resumes an interrupted plan run whose record change set has landed and the run goes on past its stopped step; in the same suite rhiza run resume still exits non-zero naming the state and rhiza run abandon for a plan run whose change set is withdrawn and for a build run whose change set has landed (run-resume-unaccepted.txtar keeps passing), each journal unchanged

## statement

rhiza run resume continues an unfinished plan-workflow run whose change set has landed (a plan run lands its record during the interview), while it still refuses a plan run whose change set is rejected, withdrawn or superseded, and any other workflow's run whose change set is no longer in flight, except a non-plan run on a landed change set that reached its land step (run-resume-landed-past-land); a non-plan run on a landed change set that never reached land is refused

## Links

### Out

- depends-on [plan-run-stops-open](plan-run-stops-open.md)
- decided-by [plan-resume-landed-record-milestone-6](../../../../../../decision/plan-resume-landed-record-milestone-6.md)

### In

- relates-to [daemon-resumes-landed-plan-run](../../../../orchestration/daemon/daemon-resumes-landed-plan-run.md) — Same exemption as rhiza run resume: plan workflow, landed state only
- verifies [plan-run-resumes-landed-record-verified](../../../../../../evidence/plan-run-resumes-landed-record-verified.md)
- verifies [plan-run-resumes-landed-record-qa-1](../../../../../../evidence/plan-run-resumes-landed-record-qa-1.md)
- verifies [plan-run-resumes-landed-record-verified-2](../../../../../../evidence/plan-run-resumes-landed-record-verified-2.md)
- verifies [pr-216-fix](../../../../../../evidence/pr-216-fix.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
