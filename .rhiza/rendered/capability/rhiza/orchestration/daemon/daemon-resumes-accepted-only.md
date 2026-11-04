# daemon-resumes-accepted-only

capability · `01M62RH3705M0GA17WVVXKBGM8` · parent [daemon](../daemon.md)

## acceptance

A test resumes an interrupted run on an accepted change set, and one on a landed change set whose non-plan run reached its land step; it leaves interrupted runs unresumed on landed change sets that never reached land, and on rejected, withdrawn, proposed and draft change sets

## statement

The daemon resumes an interrupted run only while its change set is accepted, or, for a non-plan run that reached its land step, landed (daemon-resumes-landed-past-land)

## Links

### Out

- depends-on [daemon-resume-once](daemon-resume-once.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)

### In

- blocks [review-01m6j962-r3-f3](../../../../question/review-01m6j962-r3-f3.md)
- verifies [daemon-resumes-accepted-only-verified](../../../../evidence/daemon-resumes-accepted-only-verified.md)
- verifies [daemon-resumes-accepted-only-qa-1](../../../../evidence/daemon-resumes-accepted-only-qa-1.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
