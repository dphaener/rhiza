# ship-change-workflow-tests

evidence · `01M5WHKBTBAK90JVBNSPZRYSAG`

## result

make check exit 0 at 1d113a9

## statement

The binary carries the ship-change workflow, and run start with no --workflow runs it and the journal names it builtin:ship-change.

## status

verified

## tests

TestScripts/run-ship-change; TestScripts/run-ship-change-fail; internal/workflow TestBuiltinsParse; internal/api TestWorkflowFile, TestResolveWorkflow

## Links

### Out

- verifies [ship-change-workflow](../capability/rhiza/workflows/built-in-workflows/ship-change-workflow.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
