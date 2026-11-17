# plan-workflow-applies-rulings

capability · `01M7MPK702EK9NJWH7V1HVYV4B` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview' passes with the step list, prompts, transitions and caps updated. go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous-left' passes with its fake reviser, on the apply-rulings prompt, linking each finding's node decided-by the arbiter's Decision: the run succeeds with rulings visit 1 failing, apply-rulings visit 1 passing, rulings visit 2 passing, then gate. go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-workflow|plan-workflow-pr|plan-workflow-autonomous|plan-review-visit-round|plan-milestone-backlog)$' passes, each script's 'check (visit N): pass, then gate' assertion updated to 'check (visit N): pass, then rulings' followed by 'rulings (visit 1): pass, then gate' (a clean review goes through rulings to the gate). Failure case: a new script plan-workflow-autonomous-unapplied.txtar, the same run with a reviser that writes nothing on the apply-rulings prompt, ends the run failed, no gate step was entered, the draft change set is not accepted, and run status names the unapplied Decision.

## files

internal/workflow/builtin/plan.toml, internal/workflow/builtin/plan-milestone.toml, internal/workflow/workflow_test.go, cmd/rhiza/testdata/script/plan-workflow-autonomous-left.txtar, cmd/rhiza/testdata/script/plan-workflow-autonomous-unapplied.txtar, cmd/rhiza/testdata/script/plan-workflow.txtar, cmd/rhiza/testdata/script/plan-workflow-pr.txtar, cmd/rhiza/testdata/script/plan-workflow-autonomous.txtar, cmd/rhiza/testdata/script/plan-review-visit-round.txtar, cmd/rhiza/testdata/script/plan-milestone-backlog.txtar

## must_not

Change the review, critic, check or revise steps or their caps; add a review round after the rulings; let the arbiter write anything but its Decision; change the ship steps after the gate.

## statement

The built-in plan and plan-milestone workflows pass check to a rulings command step (run = rhiza run review-check --rulings, max_visits 2, on.pass gate, on.fail apply-rulings) and add an apply-rulings agent step (prompt '/rhiza-plan-reviser Apply the arbiter's rulings on the plan of change set {change}: follow rhiza guide plan-reviser.', until true, until_when done, max_visits 1, on.pass rulings), so a ruling still unapplied after one apply visit ends the run failed before the gate.

## Links

### Out

- depends-on [plan-review-check-rulings](plan-review-check-rulings.md)
- depends-on [review-question-in-node-change-set](../../../../review-question-in-node-change-set.md)
- decided-by [plan-rulings-applied-before-gate](../../../../../decision/plan-rulings-applied-before-gate.md)
- decided-by [review-ruling-in-node-change-set](../../../../../decision/review-ruling-in-node-change-set.md)

### In

- blocks [ruling-on-serving-draft-node-unappliable](../../../../../question/ruling-on-serving-draft-node-unappliable.md)
- verifies [plan-workflow-applies-rulings-verified](../../../../../evidence/plan-workflow-applies-rulings-verified.md)
- verifies [plan-workflow-applies-rulings-qa-1](../../../../../evidence/plan-workflow-applies-rulings-qa-1.md)
- verifies [plan-workflow-applies-rulings-qa-2](../../../../../evidence/plan-workflow-applies-rulings-qa-2.md)
- verifies [plan-workflow-applies-rulings-qa-3](../../../../../evidence/plan-workflow-applies-rulings-qa-3.md)
- verifies [plan-workflow-applies-rulings-qa-4](../../../../../evidence/plan-workflow-applies-rulings-qa-4.md)
- verifies [plan-workflow-applies-rulings-qa-5](../../../../../evidence/plan-workflow-applies-rulings-qa-5.md)
- verifies [plan-workflow-applies-rulings-qa-6](../../../../../evidence/plan-workflow-applies-rulings-qa-6.md)
- verifies [plan-workflow-applies-rulings-qa-7](../../../../../evidence/plan-workflow-applies-rulings-qa-7.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
