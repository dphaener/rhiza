# arbiter-sessions

capability · `01M5WRHF3WYWZGR13RCYDS6NXY` · parent [workflows](../workflows.md)

## acceptance

A testscript runs a workflow whose fake agent raises a question: autonomous, a fake arbiter is launched with the question's ID and its Decision answers the question; normal, the run is blocked on the owner and no arbiter is launched.

## statement

In an autonomous run, the engine calls a fresh arbiter session on a session's question, told only the question's ID; it decides from slice and why and records its Decision as (assumed, would have escalated) with its reasoning. In a normal run questions go to the owner and no arbiter runs.

## Links

### Out

- decided-by [arbiter](../../../decision/arbiter.md)

### In

- verifies [pr-39-autonomous-mode-and-the](../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
