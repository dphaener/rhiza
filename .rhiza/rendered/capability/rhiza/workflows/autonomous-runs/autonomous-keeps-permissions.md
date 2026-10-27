# autonomous-keeps-permissions

capability · `01M5WRHDW5DG93SS5KHGSB709X` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

A testscript runs the same workflow with and without --autonomous on a fake adapter and finds the same --permission-mode in both launches; run start --autonomous --permission-mode bypassPermissions exits non-zero with the trust = "dedicated" error and no journal is written.

## statement

Autonomous mode never changes an agent's permissions: its agent and arbiter sessions launch with the permission mode they would have without it, and a permission mode that bypasses the agent's checks is refused with --autonomous exactly as without it, until a machine has trust = "dedicated".

## Links

### Out

- depends-on [run-start-autonomous](run-start-autonomous.md)
- decided-by [autonomous-mode](../../../../decision/autonomous-mode.md)

### In

- verifies [autonomous-keeps-permissions-verified](../../../../evidence/autonomous-keeps-permissions-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
