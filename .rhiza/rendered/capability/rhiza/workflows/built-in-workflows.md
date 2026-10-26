# built-in-workflows

capability · `01M5SWT9EFKV8YB9VXVTRK4755` · parent [workflows](../workflows.md)

## statement

The binary ships the plan, build, review and ship workflows.

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
