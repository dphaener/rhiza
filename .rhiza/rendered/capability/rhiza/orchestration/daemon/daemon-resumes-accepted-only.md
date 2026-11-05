# daemon-resumes-accepted-only

capability · `01M62RH3705M0GA17WVVXKBGM8` · parent [daemon](../daemon.md)

## acceptance

A test resumes an interrupted run on an accepted change set, an interrupted plan run on a landed change set, and an interrupted non-plan run on a landed change set whose journal reached land, and leaves an interrupted non-plan run on a landed change set that never reached land and runs on rejected, withdrawn, proposed and draft change sets unresumed

## statement

The daemon resumes an interrupted run only while its change set is accepted, except two runs on a landed change set: a plan-workflow run, which it resumes as in flight (daemon-resumes-landed-plan-run), and a non-plan run whose journal reached its land step (daemon-resumes-landed-past-land)

## Links

### Out

- depends-on [daemon-resume-once](daemon-resume-once.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)

### In

- blocks [review-01m6j962-r3-f3](../../../../question/review-01m6j962-r3-f3.md)
- blocks [review-01m6jzdh-r3-f2](../../../../question/review-01m6jzdh-r3-f2.md)
- verifies [daemon-resumes-accepted-only-verified](../../../../evidence/daemon-resumes-accepted-only-verified.md)
- verifies [daemon-resumes-accepted-only-qa-1](../../../../evidence/daemon-resumes-accepted-only-qa-1.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
