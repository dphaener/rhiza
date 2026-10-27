# build-workflow-file

capability · `01M5YX3JAF2HPFAJBN420H5JM7` · parent [build-workflow](../build-workflow.md)

## acceptance

A testscript with fake agents runs `rhiza run start build` on a change set of two leaves, one depending on the other, and finds both built and verified in order, build-check passed, the run done and the change set still accepted; `describe-golden` is regenerated.

## statement

Rhiza ships a built-in `build` workflow: a leaves step (builder models sonnet, opus, opus; verifier opus; split auto) followed by `rhiza run build-check`, ending done without landing the change set.

## Links

### Out

- depends-on [build-leaf-order](build-leaf-order.md)
- depends-on [build-retries](build-retries.md)
- depends-on [build-leaf-fails](build-leaf-fails.md)
- depends-on [build-panes](build-panes.md)
- depends-on [build-check](build-check.md)
- depends-on [build-guide](build-guide.md)
- depends-on [build-evidence-commit](build-evidence-commit.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [four-workflows](../../../../../decision/four-workflows.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- verifies [build-workflow-file-verified](../../../../../evidence/build-workflow-file-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
