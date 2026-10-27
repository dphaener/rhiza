# build-verifier-checkout

capability · `01M5YX3DGKC5HY4KRTCXA5450D` · parent [build-workflow](../build-workflow.md)

## acceptance

A test finds the verifier session's working directory a detached worktree at the leaf's commit, gone after the session ends whatever its outcome, the Evidence entry it wrote in the build worktree's .rhiza/log, and the builder's uncommitted edits in the build worktree untouched.

## statement

Each verifier works in a throwaway detached worktree at its leaf's commit, removed once its session ends, while its rhiza writes go into the build worktree's log.

## Links

### Out

- depends-on [build-leaf-commit](build-leaf-commit.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- decided-by [one-worktree-first](../../../../../decision/one-worktree-first.md)
- decided-by [answer-ab54d544-2](../../../../../decision/answer-ab54d544-2.md)
- decided-by [build-ruling-2](../../../../../decision/build-ruling-2.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-evidence-commit](build-evidence-commit.md)
- verifies [build-verifier-checkout-verified](../../../../../evidence/build-verifier-checkout-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
