# build-panes

capability · `01M5YX3GT71M5KWXZYCZA684DF` · parent [build-workflow](../build-workflow.md)

## acceptance

Tests with a fake rex find the split `right` for a landscape pane size, `below` for a portrait one, the override honored, later verifiers splitting the verifier half, a passed session's pane closed, and a failed one's left open until the run ends.

## statement

A leaves step opens one tab for the change set with the builder in its first pane and each verifier in a pane split along the longer side in pixels (the builder keeping half), unless the step's `split` fixes right or below; a pane closes when its session passes, and a failed or blocked session's pane stays until the run ends.

## Links

### Out

- depends-on [build-leaves-step](build-leaves-step.md)
- decided-by [full-agent-sessions](../../../../../decision/full-agent-sessions.md)
- decided-by [answer-ab54d544-8](../../../../../decision/answer-ab54d544-8.md)
- decided-by [answer-ab54d544-9](../../../../../decision/answer-ab54d544-9.md)
- decided-by [build-ruling-3](../../../../../decision/build-ruling-3.md)
- decided-by [build-ruling-9](../../../../../decision/build-ruling-9.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-panes-verified](../../../../../evidence/build-panes-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
