# inbox-landed-plan-run-in-flight

capability · `01M6JZNSXJ5QSSEAJF7QWJ62XT` · parent [inbox](../inbox.md)

## acceptance

A testscript with plan runs on a landed record shows no skipped item for an interrupted one, a gave-up item naming rhiza run resume for one interrupted after a daemon resume, and a failed item naming rhiza run resume --retry for a failed one; in the same suite a build run interrupted on a landed change set whose journal never reached land is still skipped with rhiza run abandon and a plan run on a withdrawn change set is still skipped, and the unaccepted-run-in-inbox and inbox-omits-resumable-landed-run tests keep passing

## statement

rhiza inbox counts a plan-workflow run whose change set has landed as in flight: it never lists it as skipped with rhiza run abandon, lists it as gave-up with rhiza run resume <run> when interrupted again after the daemon's one resume, and as failed with rhiza run resume --retry <run> when it failed, while plan runs on rejected, withdrawn or superseded change sets and other workflows' runs on terminal change sets stay skipped or hidden as before, except a non-plan run on a landed change set whose journal reached its land step, which is neither skipped nor hidden (inbox-omits-resumable-landed-run)

## Links

### Out

- relates-to [unaccepted-run-in-inbox](unaccepted-run-in-inbox.md) — Exempts a plan run on a landed record from this skipped item
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)

### In

- verifies [inbox-landed-plan-run-in-flight-verified](../../../../evidence/inbox-landed-plan-run-in-flight-verified.md)
- verifies [inbox-landed-plan-run-in-flight-qa-1](../../../../evidence/inbox-landed-plan-run-in-flight-qa-1.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
