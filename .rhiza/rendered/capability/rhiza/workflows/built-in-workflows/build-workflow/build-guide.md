# build-guide

capability · `01M5YX3HRZCS9HFYWGDSC5TGS8` · parent [build-workflow](../build-workflow.md)

## acceptance

`rhiza guide ls` lists the build topic; a testscript finds the trailer, the outside-work rule, and the Evidence-or-findings rule in it; `rhiza guide install` writes it into the plugin.

## statement

`rhiza guide build` tells builders and verifiers their part: where to start, the Leaf: trailer, needed outside work as a question and noticed-only work as a discovered-from draft change set never built in the run, and the verifier's Evidence or findings; guide install installs it.

## Links

### Out

- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-guide-verified](../../../../../evidence/build-guide-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
