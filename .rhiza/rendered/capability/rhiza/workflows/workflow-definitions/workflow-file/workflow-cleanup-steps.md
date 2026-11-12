# workflow-cleanup-steps

capability · `01M7642M6E55R7WVH65XV9BKET` · parent [workflow-file](../workflow-file.md)

## acceptance

A testscript runs a workflow whose step fails and shows its cleanup step ran and the run still failed, and a workflow whose cleanup step exits 1 still ending done with the failure in the journal; internal/workflow parse tests cover the new section.

## statement

A workflow file may declare cleanup steps, command steps that run once, in order, when a run of the workflow ends done, failed or abandoned, in the run's worktree with the run's environment; a failing cleanup step is journalled and logged and never changes the run's result; the built-in workflows declare none.

## Links

### In

- verifies [workflow-cleanup-steps-verified](../../../../../evidence/workflow-cleanup-steps-verified.md)
- verifies [workflow-cleanup-steps-qa-1](../../../../../evidence/workflow-cleanup-steps-qa-1.md)
- verifies [workflow-cleanup-steps-verified-2](../../../../../evidence/workflow-cleanup-steps-verified-2.md)
- verifies [pr-468-workflow-cleanup-phase](../../../../../evidence/pr-468-workflow-cleanup-phase.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.
