# build-check

capability · `01M5YX3HBR6XX5W3QXW7PNB5TS` · parent [build-workflow](../build-workflow.md)

## acceptance

Testscripts find build-check passing on a change set whose leaves each have Evidence by a verifier actor, and failing, naming the leaf, for a leaf with no Evidence and for one whose Evidence's actor wrote a commit with its Leaf: trailer.

## statement

`rhiza run build-check` passes when every leaf of the run's change set has live Evidence whose actor wrote none of the leaf's commits, and otherwise fails naming each leaf that lacks it.

## Links

### Out

- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- decided-by [answer-ab54d544-11](../../../../../decision/answer-ab54d544-11.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-check-verified](../../../../../evidence/build-check-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
