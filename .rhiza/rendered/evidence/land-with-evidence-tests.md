# land-with-evidence-tests

evidence · `01M5WHKD7Y22PB4S8Y0CTJPBXH`

## result

make check exit 0 at 84d07c4

## statement

The change set is landed on the PR branch with the PR and its CI result written as Evidence that verifies each capability, resume-safe and refusing to land onto main.

## status

verified

## tests

TestScripts/run-ship-change; internal/api TestShipLand, TestShipLandAfterAgentLanded, TestShipLandRefusals, TestShipEvidenceHandle, TestShipEvidenceFields

## Links

### Out

- verifies [land-with-evidence](../capability/rhiza/orchestration/shipping/land-with-evidence.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
