# build-leaf-order

capability · `01M5YX3C91E50NQQES3SAZGBDB` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test with a fake runner over leaves A, B and C, C depending on A, finds B's builder started while A's verifier is still running, and C's builder not started before A's Evidence exists.

## statement

In a leaves step builders run one at a time in the change set's leaf order, a leaf starting only once every leaf it depends-on has Evidence, while a built leaf's verifier runs alongside the next builder.

## Links

### Out

- depends-on [build-leaves-step](build-leaves-step.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [one-worktree-first](../../../../../decision/one-worktree-first.md)
- decided-by [build-ruling-1](../../../../../decision/build-ruling-1.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-leaf-order-verified](../../../../../evidence/build-leaf-order-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
