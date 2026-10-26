# ship-change-workflow

capability · `01M5TYKP0BBSKRGKDXXCBVWT65` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

A testscript starts a run without naming a workflow and its journal names ship-change.

## statement

The binary ships the ship-change workflow, and rhiza run start <change set> runs it by default.

## Links

### Out

- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)
- decided-by [built-in-workflows-by-name](../../../../decision/built-in-workflows-by-name.md)

### In

- verifies [ship-change-workflow-tests](../../../../evidence/ship-change-workflow-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
