# land-with-evidence

capability · `01M5TYKQJ8YK6B2WW3AX1Y3C6X` · parent [shipping](../shipping.md)

## acceptance

After a fake run the change set is landed and an Evidence node names the PR and the checks.

## statement

Ship-change lands the change set on the PR branch once its CI passes, before the merge, and records the PR and its CI result as Evidence that verifies it.

## Links

### Out

- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)
- decided-by [one-change-set-one-pr](../../../../decision/one-change-set-one-pr.md)
- decided-by [ship-lands-before-merge](../../../../decision/ship-lands-before-merge.md)

### In

- verifies [land-with-evidence-tests](../../../../evidence/land-with-evidence-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
