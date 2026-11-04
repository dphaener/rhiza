# unaccepted-run-in-inbox

capability · `01M62RH3PBJ12H0CYSMCMV7KAS` · parent [inbox](../inbox.md)

## acceptance

A test with an interrupted run on a landed change set whose journal never reached land lists it in rhiza inbox with the state and the abandon command, and the daemon does not resume it

## statement

An interrupted run on a change set no longer accepted waits in the inbox, naming the change set's state and rhiza run abandon as the next step, except a non-plan run on a landed change set whose journal reached its land step, which the daemon resumes (inbox-omits-resumable-landed-run, daemon-resumes-landed-past-land)

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)

### In

- blocks [review-01m6j962-r3-f1](../../../../question/review-01m6j962-r3-f1.md)
- verifies [unaccepted-run-in-inbox-verified](../../../../evidence/unaccepted-run-in-inbox-verified.md)
- verifies [unaccepted-run-in-inbox-qa-1](../../../../evidence/unaccepted-run-in-inbox-qa-1.md)
- verifies [unaccepted-run-in-inbox-verified-2](../../../../evidence/unaccepted-run-in-inbox-verified-2.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
