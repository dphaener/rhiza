# handback-clears-done-mark

decision · `01M7MD742M9ZVWG0J0HSJA2SB6`

## statement

A hand-back that messages a live session clears that session's done mark (Status.Done) before the message is typed, as resuming an ended session already clears its whole status; the asked mark and the hook, rex and start reports are left as they are, so the visit counts as done only on a fresh rhiza session done

## wording

(derived from code: the resume path already clears the status; the live path is the gap the bug names)

## Links

### Out

- derived-from [rejoined-conductor-done-mark-plan](../source/rejoined-conductor-done-mark-plan.md)

### In

- decided-by [handback-clears-live-done-mark](../capability/rhiza/workflows/workflow-runs/run-handback/handback-clears-live-done-mark.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
