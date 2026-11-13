# build-guide

capability · `01M5YX3HRZCS9HFYWGDSC5TGS8` · parent [build-workflow](../build-workflow.md)

## acceptance

`rhiza guide ls` lists the build, builder, verifier and build-retry topics; go test -trimpath ./internal/guide/... finds the trailer, the outside-work rule and the Evidence-or-findings rule in the role topics and the three names in the build index; `rhiza guide install` writes all four as skills.

## statement

`rhiza guide build` is the build index: it names the builder, verifier and build-retry topics and says which a builder, a verifier and a retrying builder load; the role topics carry where to start, the Leaf: trailer, needed outside work as a question and noticed-only work as a discovered-from draft change set never built in the run, and the verifier's Evidence or findings; guide install installs them all.

## Links

### Out

- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-guide-verified](../../../../../evidence/build-guide-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)
- verifies [pr-521-build-role-topics](../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
