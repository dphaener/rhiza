# close-out

capability · `01M5TYKQ3V9NFPG905CN1RWR3B` · parent [shipping](../shipping.md)

## acceptance

A testscript with fakes ends with no branch, no worktree and a green main.

## statement

After the merge the branch is deleted locally and on GitHub, the worktree is removed, and CI on main is waited on until green.

## Links

### Out

- decided-by [task-close-out](../../../../decision/task-close-out.md)

### In

- verifies [close-out-tests](../../../../evidence/close-out-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
