# close-out-tests

evidence · `01M5WHKCS2YZ3B5EKX96BWSJ3X` · actor [claude-code-4c6d8166](../actor/claude-code-4c6d8166.md)

## result

make check exit 0 at 1d113a9

## statement

After the merge the worktree, the session's Rex tab, the local branch and the origin branch are gone, main is pulled, and CI on main is waited on until green; a tab Rex no longer has counts as closed.

## status

verified

## tests

TestScripts/run-ship-change; internal/api TestShipCloseOut, TestShipCloseOutDirtyWorktree, TestShipCloseOutClosesTabs

## Links

### Out

- verifies [close-out](../capability/rhiza/orchestration/shipping/close-out.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
