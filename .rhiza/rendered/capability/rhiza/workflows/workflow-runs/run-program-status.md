# run-program-status

capability · `01M5WT1XEJY38RKEE9K473K3V1` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A test captures the bytes a run writes with a fake Rex through a build of three leaves and finds, in order, root working with the step, child records for each leaf with builder then verifier, progress rising to 100, a blocked child while a fake session waits, and root done; a failing run ends with root error.

## statement

While a run goes, its process keeps Rex's program status current through OSC 7501: a root record with the current step in its message and, in the build workflow, progress as the share of leaves with Evidence; a child record per live leaf (hierarchical ID per the spec) saying which leaf, builder or verifier, and working or blocked; and a final root state of done when the run succeeds or error when it fails.

## Links

### Out

- depends-on [build-leaves-step](../built-in-workflows/build-workflow/build-leaves-step.md)
- decided-by [run-reports-progress](../../../../decision/run-reports-progress.md)

### In

- verifies [run-program-status-verified](../../../../evidence/run-program-status-verified.md)
- verifies [pr-44-the-build-workflow](../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
