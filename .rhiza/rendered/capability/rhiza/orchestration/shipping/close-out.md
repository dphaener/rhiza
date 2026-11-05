# close-out

capability · `01M5TYKQ3V9NFPG905CN1RWR3B` · parent [shipping](../shipping.md)

## acceptance

A testscript with fakes ends with no branch, no worktree, the session's Rex tab closed and main confirmed by the merge commit's tree.

## statement

After the merge the branch is deleted locally and on GitHub, the worktree is removed, and main is confirmed by checking that the merge commit's tree on main matches the shipped branch head's tree.

## Links

### Out

- decided-by [task-close-out](../../../../decision/task-close-out.md)
- decided-by [close-out-closes-tab](../../../../decision/close-out-closes-tab.md)
- decided-by [close-out-fast-forwards-main](../../../../decision/close-out-fast-forwards-main.md)

### In

- verifies [close-out-tests](../../../../evidence/close-out-tests.md)
- verifies [pr-256-update-github-runner-to](../../../../evidence/pr-256-update-github-runner-to.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
