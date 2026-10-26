# shipit-agent-step

capability · `01M5TYKPAXPA4NBWNPH5WENFNK` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

With a fake adapter, the step's prompt names shipit and the slice command.

## statement

Ship-change's agent step runs shipit in the change set's worktree, starting from rhiza slice --change.

## Links

### Out

- decided-by [shipit-as-agent-step](../../../../decision/shipit-as-agent-step.md)

### In

- verifies [shipit-agent-step-tests](../../../../evidence/shipit-agent-step-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
