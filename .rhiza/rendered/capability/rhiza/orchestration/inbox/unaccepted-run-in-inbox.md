# unaccepted-run-in-inbox

capability · `01M62RH3PBJ12H0CYSMCMV7KAS` · parent [inbox](../inbox.md)

## acceptance

A test with an interrupted build run on a landed change set whose journal never reached land lists it in rhiza inbox with the state and the abandon command, and the daemon does not resume it; an interrupted plan run on a landed change set and an interrupted non-plan run on a landed change set whose journal reached land are not listed as skipped

## statement

An interrupted run on a change set no longer accepted waits in the inbox, naming the change set's state and rhiza run abandon as the next step, except two runs on a landed change set, which the daemon resumes: a plan-workflow run, which counts as in flight (inbox-landed-plan-run-in-flight), and a non-plan run whose journal reached its land step (inbox-omits-resumable-landed-run)

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- decided-by [skipped-run-in-inbox](../../../../decision/skipped-run-in-inbox.md)
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)

### In

- relates-to [inbox-landed-plan-run-in-flight](inbox-landed-plan-run-in-flight.md) — Exempts a plan run on a landed record from this skipped item
- blocks [review-01m6j962-r3-f1](../../../../question/review-01m6j962-r3-f1.md)
- verifies [unaccepted-run-in-inbox-verified](../../../../evidence/unaccepted-run-in-inbox-verified.md)
- verifies [unaccepted-run-in-inbox-qa-1](../../../../evidence/unaccepted-run-in-inbox-qa-1.md)
- verifies [unaccepted-run-in-inbox-verified-2](../../../../evidence/unaccepted-run-in-inbox-verified-2.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
