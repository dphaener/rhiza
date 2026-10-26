# repo-ship-command-tests

evidence · `01M5WHKD0PWADHFWMYHE6KW2D7`

## result

make check exit 0 at 84d07c4

## statement

The repo's .rhiza/workflows/ship.sh runs once, after main is green, in the main checkout; without it the command reports there is nothing to ship.

## status

verified

## tests

TestScripts/run-ship-change; TestScripts/run-ship-change-fail; internal/api TestShipCommand

## Links

### Out

- verifies [repo-ship-command](../capability/rhiza/orchestration/shipping/repo-ship-command.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
