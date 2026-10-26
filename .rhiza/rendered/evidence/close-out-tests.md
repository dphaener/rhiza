# close-out-tests

evidence · `01M5WHKCS2YZ3B5EKX96BWSJ3X`

## result

make check exit 0 at 84d07c4

## statement

After the merge the worktree, local branch and origin branch are gone, main is pulled, and CI on main is waited on until green.

## status

verified

## tests

TestScripts/run-ship-change; internal/api TestShipCloseOut, TestShipCloseOutDirtyWorktree

## Links

### Out

- verifies [close-out](../capability/rhiza/orchestration/shipping/close-out.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
