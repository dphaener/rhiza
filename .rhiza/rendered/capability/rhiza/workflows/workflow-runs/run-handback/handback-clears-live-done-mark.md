# handback-clears-live-done-mark

capability · `01M7MD744F3AQHKM0BZJ7E03FN` · parent [run-handback](../run-handback.md)

## acceptance

go test -trimpath ./internal/api -run 'TestHandbackDeliver' passes with a new subtest: a live idle session whose record carries a done mark is handed back; afterwards the record has no done mark, the hand-back message is written and typed, and stepRunner.AgentStatus reports SaidDone false; then rhiza session done's mark set again makes AgentStatus report SaidDone true. (failure case) a hand-back to the same done-marked live session through a step with no handback template fails naming the step and leaves the done mark in place, with no message written.

## files

internal/api/run.go (Handback, messageSession), internal/api/run_handback_test.go

## must_not

Clear the asked mark or the hook, rex, started or prompted reports; change the resume path for an ended session; change rhiza session done, session message or the engine in internal/run; clear the done mark when the hand-back fails before any message is sent

## statement

When a run hands back to a live session that said done on an earlier visit, stepRunner.Handback clears the record's done mark before messaging it, so AgentStatus reports it not said done and a step waiting for done (ship-change's conduct) passes only after the session says done again

## Links

### Out

- decided-by [handback-clears-done-mark](../../../../../decision/handback-clears-done-mark.md)

### In

- verifies [handback-clears-live-done-mark-verified](../../../../../evidence/handback-clears-live-done-mark-verified.md)
- verifies [handback-clears-live-done-mark-qa-1](../../../../../evidence/handback-clears-live-done-mark-qa-1.md)
- verifies [handback-clears-live-done-mark-qa-2](../../../../../evidence/handback-clears-live-done-mark-qa-2.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
