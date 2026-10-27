# build-workflow

capability · `01M5YX3AWPBAH862Q9ZEHM84NW` · parent [built-in-workflows](../built-in-workflows.md)

## statement

The build workflow builds a change set's leaves: a builder session and then a fresh verifier session per leaf, in panes of the change set's tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.

## Links

### Out

- decided-by [build-workflow-shape](../../../../decision/build-workflow-shape.md)
- decided-by [build-ruling-5](../../../../decision/build-ruling-5.md)
- decided-by [build-ruling-6](../../../../decision/build-ruling-6.md)
- decided-by [build-ruling-8](../../../../decision/build-ruling-8.md)
- decided-by [build-signoff](../../../../decision/build-signoff.md)
- derived-from [build-workflow-plan](../../../../source/build-workflow-plan.md)

### In

- verifies [pr-44-the-build-workflow](../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
