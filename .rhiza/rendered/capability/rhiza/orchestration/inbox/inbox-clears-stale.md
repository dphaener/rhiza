# inbox-clears-stale

capability · `01M62S5BN8M65JXT1672T8GVEF` · parent [inbox](../inbox.md)

## acceptance

A testscript with fixture journals and change sets: a failed run and a gave-up run each on a landed, a withdrawn, a rejected and a superseded change set, and one of each on an accepted change set, lists exactly the accepted-change-set items in text and --json; a failed and a gave-up plan run on a landed change set stay listed, and so do a failed and a gave-up non-plan run on a landed change set whose journal reached land; a run on a change set the log does not know stays listed; a question and a sign-off on a landed change set stay listed

## statement

rhiza inbox leaves out a failed run or a run the daemon gave up on whose change set has landed, been withdrawn, rejected or superseded, so every item it lists names a next command that can succeed on its change set; two runs on a landed change set are not left out: a plan-workflow run, since rhiza run resume works on it (inbox-landed-plan-run-in-flight), and a non-plan run whose journal reached its land step (inbox-omits-resumable-landed-run); questions and sign-offs stay listed

## Links

### Out

- decided-by [answer-29b4ed69-1](../../../../decision/answer-29b4ed69-1.md)
- decided-by [answer-29b4ed69-2](../../../../decision/answer-29b4ed69-2.md)
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)
- derived-from [inbox-stale-interview](../../../../source/inbox-stale-interview.md)

### In

- blocks [review-01m6jzdh-r3-f4](../../../../question/review-01m6jzdh-r3-f4.md)
- verifies [inbox-clears-stale-verified](../../../../evidence/inbox-clears-stale-verified.md)
- verifies [inbox-clears-stale-qa-1](../../../../evidence/inbox-clears-stale-qa-1.md)
- verifies [pr-66-inbox-clears-stale-failures](../../../../evidence/pr-66-inbox-clears-stale-failures.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
