# build-retries

capability · `01M5YX3F9R6BSD6XKKTVQXYSDZ` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test whose verifier fails twice then passes finds three builder sessions on the ladder's models, each prompt naming its attempt, and the leaf passed; one whose verifier fails three times finds no fourth builder and the leaf failed.

## statement

A failed attempt goes to a fresh builder session on the next model of the ladder, its prompt naming the leaf and the attempt, which reads the earlier attempts' findings from the run; after 3 failed attempts the leaf fails.

## Links

### Out

- depends-on [build-model-ladder](build-model-ladder.md)
- depends-on [build-verifier-findings](build-verifier-findings.md)
- depends-on [build-leaf-commit](build-leaf-commit.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [answer-ab54d544-3](../../../../../decision/answer-ab54d544-3.md)
- decided-by [answer-ab54d544-5](../../../../../decision/answer-ab54d544-5.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-leaf-fails](build-leaf-fails.md)
- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-retries-verified](../../../../../evidence/build-retries-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
