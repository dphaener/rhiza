# run-resume

capability · `01M5TYGK7EFHKPB9JG3A2MH2G2` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript stops a run after step one, resumes it, and step one runs once.

## statement

`genesis run resume <run>` continues a stopped run from its journal without a daemon, never redoing a finished step; it refuses an interrupted run whose change set has landed, been rejected, withdrawn or superseded, except that a run whose change set has landed is resumed when it is a plan-workflow run (plan-run-resumes-landed-record) or reached its land step (run-resume-landed-past-land).

## Links

### Out

- decided-by [run-journal](../../../../decision/run-journal.md)

### In

- depends-on [manual-resume-refuses-unaccepted](manual-resume-refuses-unaccepted.md)
- depends-on [run-resume-landed-past-land](run-resume-landed-past-land.md)
- verifies [run-resume-tested](../../../../evidence/run-resume-tested.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-216-fix](../../../../evidence/pr-216-fix.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
