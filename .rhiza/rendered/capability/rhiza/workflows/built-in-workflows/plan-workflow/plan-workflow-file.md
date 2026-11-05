# plan-workflow-file

capability · `01M5WSGDJMV9EEVE603M22PTYT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go workflow tests: BuiltinNames lists plan and its file parses with steps interview, milestone-check, review, critic, check, revise, gate, open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out, the gate passing to open-pr and close-out ending the run; the review step's prompt names only {change} and rhiza guide plan-review, and the interview step's prompt names the rhiza-plan skill. A testscript runs rhiza run start --workflow plan with fake agents and a fake gh on a draft change set through the gate to open-pr, and the run's branch is rendered and committed: git status is clean, rhiza check passes, and no plan file is written.

## statement

The binary ships a built-in plan workflow: an interview agent step driven by the rhiza-plan skill, then a milestone-check command step (rhiza run milestone-check) that ends the run failed when a change set the plan makes serves no milestone, then rounds of a fresh review agent step told only the change set's ID, a critic agent step and rhiza run review-check, with a revise agent step on a failed check, then rhiza run gate, and after it the ship steps open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out. It writes no plan file.

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
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
