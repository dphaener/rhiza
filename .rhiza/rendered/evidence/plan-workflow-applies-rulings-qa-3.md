# plan-workflow-applies-rulings-qa-3

evidence · `01M7MTBE39WNRTZABBF0X31M72` · actor [claude-code-60928770](../actor/claude-code-60928770.md)

## result

PASS

## statement

QA round 3 at db1b1c925, scratch build (go build -trimpath, run by full path, outside PATH) driving autonomous plan runs by hand in scratch Rhiza roots with the plan-workflow-autonomous* fake rex/claude/gh: rhiza workflow describe plan and plan-milestone show check→rulings→gate with rulings on.fail apply-rulings and the apply-rulings prompt; a run whose reviser links each finding's node decided-by went rulings fail → apply-rulings → rulings pass → gate (both change sets accepted); one whose reviser writes nothing went rulings fail → apply-rulings → rulings fail and failed ('apply-rulings went over max_visits 1') with no gate step, both change sets still draft and the unapplied Decisions named. go test ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview' ok.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
