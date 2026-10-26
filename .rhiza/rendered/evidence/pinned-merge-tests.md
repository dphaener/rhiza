# pinned-merge-tests

evidence · `01M5WHKCHS79BZTVXXWHAG4DG5`

## result

make check exit 0 at 1d113a9

## statement

The merge is gh pr merge --squash --match-head-commit with the checked head, and a merge that did not happen is an error.

## status

verified

## tests

TestScripts/run-ship-change; internal/api TestShipCheckGreenThenMerge, TestShipMergeNotMerged

## Links

### Out

- verifies [pinned-merge](../capability/rhiza/orchestration/shipping/pinned-merge.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
