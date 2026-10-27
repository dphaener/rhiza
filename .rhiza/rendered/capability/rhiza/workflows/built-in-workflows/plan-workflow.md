# plan-workflow

capability · `01M5WSGC24VYV3Y9D1N4WQ39PC` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

A testscript runs rhiza run start --workflow plan on a draft change set with fake agents that interview, review with one finding the critic confirms, revise, and review clean: the run succeeds, the plan's change sets are proposed, git status in the worktree shows changes only under .rhiza/, and rhiza check passes.

## statement

rhiza run start --workflow plan runs Rhiza's plan workflow: an agent interviews the owner one topic at a time and writes a milestone, or a single change set's leaves, straight into the model, with a Source and Decisions quoting the owner's typed words, a fresh review with a critic and at most 3 rounds, and the owner's acceptance as the gate. It keeps no plan file and no checkboxes.

## Links

### Out

- decided-by [plan-workflow-shape](../../../../decision/plan-workflow-shape.md)
- decided-by [model-is-plan-and-state](../../../../decision/model-is-plan-and-state.md)
- decided-by [four-workflows](../../../../decision/four-workflows.md)

### In

- verifies [plan-workflow-verified](../../../../evidence/plan-workflow-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
