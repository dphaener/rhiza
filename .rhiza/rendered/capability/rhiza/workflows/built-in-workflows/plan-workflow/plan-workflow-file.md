# plan-workflow-file

capability · `01M5WSGDJMV9EEVE603M22PTYT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go workflow tests: BuiltinNames lists plan and its file parses with steps interview, review, critic, check, revise, gate, open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out, the gate passing to open-pr and close-out ending the run; the review step's prompt names only {change} and rhiza guide plan-review, and the interview step's prompt names the rhiza-plan skill. A testscript runs rhiza run start --workflow plan with fake agents and a fake gh on a draft change set through the gate to open-pr, and the run's branch is rendered and committed: git status is clean, rhiza check passes, and no plan file is written.

## statement

The binary ships a built-in plan workflow: an interview agent step driven by the rhiza-plan skill, then rounds of a fresh review agent step told only the change set's ID, a critic agent step and rhiza run review-check, with a revise agent step on a failed check, then rhiza run gate, and after it the ship steps open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out. It writes no plan file.

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
- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
