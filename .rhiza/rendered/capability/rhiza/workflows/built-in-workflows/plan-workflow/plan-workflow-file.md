# plan-workflow-file

capability · `01M5WSGDJMV9EEVE603M22PTYT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go workflow tests: BuiltinNames lists plan and its file parses with steps interview, milestone-check, review, critic, check, revise, gate, open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out, the gate passing to open-pr and close-out ending the run; the review step's prompt names only {change} and rhiza guide plan-review, and the interview step's prompt names the rhiza-plan skill. A testscript runs rhiza run start --workflow plan with fake agents and a fake gh on a draft change set through the gate to open-pr, and the run's branch is rendered and committed: git status is clean, rhiza check passes, and no plan file is written.

## statement

The binary ships a built-in plan workflow: an interview agent step driven by the rhiza-plan skill, then a milestone-check command step (rhiza run milestone-check) that ends the run failed only as plan-step-milestone-check sets out (a planned change set serving no milestone passes, unless the run was started with a milestone it does not serve or the repo's [schedule] backlog is "skip"), then rounds of a fresh review agent step told only the change set's ID, a critic agent step and rhiza run review-check, with a revise agent step on a failed check, then rhiza run gate, and after it the ship steps open-pr, wait-pr, pr-green, rebase, wait-rebased, rebased-green, merge and close-out. It writes no plan file.

## Links

### Out

- depends-on [plan-guide-install](plan-guide-install.md)
- depends-on [plan-review-guide](plan-review-guide.md)
- depends-on [plan-gate](plan-gate.md)
- depends-on [milestone-check-follows-policy](../../../orchestration/daemon/plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [model-is-plan-and-state](../../../../../decision/model-is-plan-and-state.md)
- decided-by [four-workflows](../../../../../decision/four-workflows.md)
- decided-by [plan-review-rounds-journalled](../../../../../decision/plan-review-rounds-journalled.md)

### In

- verifies [plan-workflow-file-verified](../../../../../evidence/plan-workflow-file-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)
- verifies [pr-358-scheduling-as-policy](../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
