# ship-command-file

decision · `01M5WHKECCYCDGPBCEH2E92FY1`

## statement

The repo's ship command is .rhiza/workflows/ship.sh, run with sh in the main checkout after main is green; without it there is nothing to ship.

## wording

the building agent's design, assumed in autonomous mode

## Links

### In

- decided-by [repo-ship-command](../capability/rhiza/orchestration/shipping/repo-ship-command.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
