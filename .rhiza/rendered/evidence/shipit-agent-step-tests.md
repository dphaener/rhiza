# shipit-agent-step-tests

evidence · `01M5WHKC1RBD6E4939DF90GQQD`

## result

make check exit 0 at 84d07c4

## statement

The agent step runs shipit in the change set's worktree from rhiza slice --change; the fake claude receives that prompt and pushes the branch.

## status

verified

## tests

TestScripts/run-ship-change; internal/workflow TestBuiltinsParse

## Links

### Out

- verifies [shipit-agent-step](../capability/rhiza/workflows/built-in-workflows/shipit-agent-step.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
