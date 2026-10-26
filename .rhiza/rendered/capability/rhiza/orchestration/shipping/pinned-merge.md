# pinned-merge

capability · `01M5TYKPWPQZ4D8ZTFYP2SWH9N` · parent [shipping](../shipping.md)

## acceptance

The fake gh sees --squash --match-head-commit <sha> with the checked SHA.

## statement

The merge is a squash pinned to the head commit whose checks passed, so a later push is never merged unchecked.

## Links

### Out

- decided-by [orchestrator-merges](../../../../decision/orchestrator-merges.md)
- decided-by [ship-state-file](../../../../decision/ship-state-file.md)

### In

- verifies [pinned-merge-tests](../../../../evidence/pinned-merge-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
