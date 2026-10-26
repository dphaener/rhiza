# close-out-fast-forwards-main

decision · `01M5WJJH2DGK6848BWH1S5GD3A`

## statement

Close-out fast-forwards the main checkout (git pull --ff-only origin main) only when it is on main with no changes to tracked files; otherwise the step fails and says why.

## wording

the building agent's design, accepted by the orchestrator on PR #30 (the owner authorized it to approve plans)

## Links

### In

- decided-by [close-out](../capability/rhiza/orchestration/shipping/close-out.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
