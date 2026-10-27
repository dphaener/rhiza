# build-leaf-commit

capability · `01M5YX3CYY6077G7P5ZTRA65M9` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test against a real git repo finds a builder that ends with no new commit, or with one naming another leaf, counted as a failed attempt with that detail, and a builder whose commit names its leaf followed by a verifier given that commit.

## statement

A builder attempt passes only when it leaves a new commit whose `Leaf:` trailer names its leaf, and a retry adds a fix commit with the same trailer; the verifier starts on that commit.

## Links

### Out

- depends-on [build-leaves-step](build-leaves-step.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [answer-ab54d544-6](../../../../../decision/answer-ab54d544-6.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-verifier-checkout](build-verifier-checkout.md)
- depends-on [build-retries](build-retries.md)
- verifies [build-leaf-commit-verified](../../../../../evidence/build-leaf-commit-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
