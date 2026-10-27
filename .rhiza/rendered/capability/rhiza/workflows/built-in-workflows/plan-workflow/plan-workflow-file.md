# plan-workflow-file

capability · `01M5WSGDJMV9EEVE603M22PTYT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go workflow tests: BuiltinNames lists plan and its file parses; the review step's prompt names only {change} and rhiza guide plan-review, and the interview step's prompt names the rhiza-plan skill. A testscript runs rhiza run start --workflow plan with fake agents on a draft change set to the gate, and git status in the worktree shows changes only under .rhiza/.

## statement

The binary ships a built-in plan workflow: an interview agent step driven by the rhiza-plan skill, then rounds of a fresh review agent step told only the change set's ID, a critic agent step and rhiza run review-check, with a revise agent step on a failed check, and finally rhiza run gate. It writes no plan file.

## Links

### Out

- depends-on [plan-guide-install](plan-guide-install.md)
- depends-on [plan-review-guide](plan-review-guide.md)
- depends-on [plan-gate](plan-gate.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [model-is-plan-and-state](../../../../../decision/model-is-plan-and-state.md)
- decided-by [four-workflows](../../../../../decision/four-workflows.md)
- decided-by [plan-review-rounds-journalled](../../../../../decision/plan-review-rounds-journalled.md)

### In

- verifies [plan-workflow-file-verified](../../../../../evidence/plan-workflow-file-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
