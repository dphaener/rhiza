# spawn-ends-on-done-mark

decision · `01M7M2DV1SVP6JQHDTYE143HB0`

## reasoning

The intent names the hang: spawn's loop waits on childOutcome, which ends only on an exit, and a child that said done sits at its prompt with no daemon to close it. The engine already counts a done mark as finished and closes the session's pane at once (readyForUntil, agent, closeSessionPane), and session.done is documented as the agent's very last action, so ending the wait on it and closing the pane follows the existing pattern; nothing hangs on the close, as in the engine.

## statement

rhiza run spawn treats a child's done mark (rhiza session done) as the end of the child, as the engine treats it for a step with until_when = "done": the wait ends with outcome pass and detail "said done", the end is journaled, and spawn then closes the child's Rex pane itself (closeSessionPane), so the child's agent exits with no daemon. A child that exited before saying done keeps today's outcome from its exit. childOutcome, shared with stepRunner.Child, reads the done mark too, so a restarted driver's reconcileSpawns records the same end. A pane close that fails leaves the outcome as it is and is only reported.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [spawn-done-plan](../source/spawn-done-plan.md)

### In

- decided-by [spawn-ends-on-done](../capability/rhiza/workflows/workflow-runs/spawn-ends-on-done.md)

## History

- `01M7M27J4MMS8X5PT4W3EXKDTM` Fix rhiza run spawn waits for the child agent's process to exit, and with no daemon nothing closes a child that already said rhiza session done, so spawn and its run hang until someone closes the child's pane
