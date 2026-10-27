# build-model-ladder

capability · `01M5YX3BVJ20HMXE6GY00WB69Y` · parent [build-workflow](../build-workflow.md)

## acceptance

Adapter tests find `--model opus` in the argv for a step naming opus and no `--model` for one naming none; an engine test finds attempts 1, 2 and 3 started on sonnet, opus and opus from `models = ["sonnet", "opus"]`.

## statement

An agent step or a leaves sub-table names models by attempt, the last reused past the end of the list, and the Claude adapter starts the session with that model.

## Links

### Out

- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [answer-ab54d544-4](../../../../../decision/answer-ab54d544-4.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-retries](build-retries.md)
- verifies [build-model-ladder-verified](../../../../../evidence/build-model-ladder-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
