# wait-green-pr-tests

evidence · `01M5WHKC9HFN1PSXE1XBJP4MYK`

## result

make check exit 0 at 84d07c4

## statement

The run waits until the PR is open and all checks pass, and fails the pr-green step on a failed check without merging.

## status

verified

## tests

TestScripts/run-ship-change; TestScripts/run-ship-change-fail; internal/api TestShipCheckNoPR, TestShipCheckPending, TestShipCheckFailed, TestShipCheckWaitsForPushedHead, TestChecksVerdict

## Links

### Out

- verifies [wait-green-pr](../capability/rhiza/orchestration/shipping/wait-green-pr.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
