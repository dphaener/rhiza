# autonomous-runs

capability · `01M5WRHD790X47F8VK7RQEHEHX` · parent [workflows](../workflows.md)

## acceptance

A testscript starts a run with --autonomous on a fake adapter whose step accepts a change set under rule autonomous: run status --json shows the run autonomous and lists that accept, and the agent's launch carries the same --permission-mode as without the flag.

## statement

A run can be started autonomous: it skips the owner's gates, recording each with rule autonomous, never bypasses an agent's permissions, and lists afterwards everything it skipped and assumed.

## Links

### Out

- decided-by [autonomous-mode](../../../decision/autonomous-mode.md)
- decided-by [two-owner-gates](../../../decision/two-owner-gates.md)

### In

- verifies [pr-39-autonomous-mode-and-the](../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
