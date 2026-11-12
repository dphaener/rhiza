# workflow-check-cleanup

capability · `01M7642M6S0FA2TZGYSZT6X311` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

internal/workflow and workflow-check testscripts cover a valid and an invalid cleanup section; describe-golden and the generated reference are regenerated.

## statement

rhiza workflow check and workflow describe cover cleanup steps: check rejects a cleanup step that is not a command step or names a next step, and describe and the generated workflow-file reference show the cleanup section.

## Links

### In

- verifies [workflow-check-cleanup-verified](../../../../evidence/workflow-check-cleanup-verified.md)
- verifies [workflow-check-cleanup-qa-1](../../../../evidence/workflow-check-cleanup-qa-1.md)
- verifies [pr-468-workflow-cleanup-phase](../../../../evidence/pr-468-workflow-cleanup-phase.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.
