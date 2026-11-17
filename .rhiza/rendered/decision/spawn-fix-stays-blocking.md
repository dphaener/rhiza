# spawn-fix-stays-blocking

decision · `01M7M2FV0RQF0FPQZ7HEWP6KV5`

## quote

yes, pass those ideas to the planner

## quote_where

the owner to the orchestrator, 2026-11-17, handing the ideas to this planner (planner-notes 01M7EAYT013); the owner's words cover only the hand-off, not this choice

## reasoning

The intent is the hang: spawn waits for an exit that a done child never makes, and nothing closes it. The orchestrator's 2026-11-17 finding (run 01M7KHWS2DRWVTKAKNSZ99TR7W: a builder done at 01:46:45Z stayed up until about 02:29Z with the daemon on) has the same cause, and spawn-ends-on-done fixes both, since spawn runs in the caller's process. A non-blocking spawn changes the conductor's protocol (rhiza-conduct, run.spawn's response, session.done's side effects), and check-in alarms add a new timer surface. Both reach past this bug's intent, so they are drafted for their own plans, where the requirement that every child stays a Rhiza session on the run's ledger can be planned in full.

## statement

This change set keeps rhiza run spawn blocking and fixes only how it ends: on the child's done mark, closing the child's pane from the spawn process, which also covers runs with the daemon on. The orchestrator's two rex-teams ideas, a spawn that returns at start with the child waking the conductor through rhiza session message, and check-in alarms on a conductor's own children, are drafted as their own change sets (01M7M2FB5BNRSDSN7HFWQSC2WN and 01M7M2FE1SMCJ681K15CBZ3R6D, discovered-from this one), not built here.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [spawn-done-plan](../source/spawn-done-plan.md)

### In

- decided-by [spawn-ends-on-done](../capability/rhiza/workflows/workflow-runs/spawn-ends-on-done.md)

## History

- `01M7M27J4MMS8X5PT4W3EXKDTM` Fix rhiza run spawn waits for the child agent's process to exit, and with no daemon nothing closes a child that already said rhiza session done, so spawn and its run hang until someone closes the child's pane
