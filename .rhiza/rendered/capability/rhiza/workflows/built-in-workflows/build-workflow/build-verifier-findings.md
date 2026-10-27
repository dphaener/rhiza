# build-verifier-findings

capability · `01M5YX3ER7D48Q8364S1JXKAW7` · parent [build-workflow](../build-workflow.md)

## acceptance

Testscripts find `rhiza run findings` with a leaf journal findings against that leaf's current attempt and `rhiza run status` show them; an engine test finds an attempt failed when its verifier ends with neither, and when the leaf's Evidence names the builder's actor.

## statement

A verifier ends a leaf's attempt with Evidence written as its own actor, or with findings it journals against the leaf and attempt; a verifier that ends with neither fails the attempt, and Evidence written by any other actor does not pass it.

## Links

### Out

- depends-on [build-leaves-step](build-leaves-step.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- decided-by [answer-ab54d544-5](../../../../../decision/answer-ab54d544-5.md)
- decided-by [build-ruling-4](../../../../../decision/build-ruling-4.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-retries](build-retries.md)
- verifies [build-verifier-findings-verified](../../../../../evidence/build-verifier-findings-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
