# wait-green-pr

capability · `01M5TYKPJ0FNX289H8PNJRM1TT` · parent [shipping](../shipping.md)

## acceptance

A testscript with a fake gh reaches pass and fail.

## statement

Ship-change waits until the change set's PR is open and every CI check has passed, and fails the step on a failed check.

## Links

### Out

- decided-by [orchestrator-merges](../../../../decision/orchestrator-merges.md)
- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)

### In

- verifies [wait-green-pr-tests](../../../../evidence/wait-green-pr-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
