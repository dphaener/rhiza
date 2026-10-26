# ship-lands-before-merge

decision · `01M5WHKDW15YK8DX5DY3NTMP3B`

## statement

Ship-change lands the change set on the PR branch after its CI passes, pushes that, waits for CI again and merges pinned to that head, so main gets the model and code in one reviewed squash.

## wording

the building agent's design, assumed in autonomous mode

## Links

### In

- decided-by [land-with-evidence](../capability/rhiza/orchestration/shipping/land-with-evidence.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
