# plan-reviser-applies-rulings-qa-5

evidence · `01M7MV603GJ99RPAS1V9VMVHQF` · actor [claude-code-602c1c12](../actor/claude-code-602c1c12.md)

## result

PASS

## statement

QA cycle 5, scratch binary go build -trimpath -o $TMPDIR/qa5/rhiza ./cmd/rhiza, scratch autonomous plan run (rhiza run start --workflow plan --autonomous, fake agents from plan-workflow-autonomous-unapplied.txtar with a reviser paused at apply-rulings): run round 3 left two arbitrated questions, one on serving-draft leaf export-csv and one on landed decision keep-flat. From the reviser's environment (RHIZA_CHANGE = the run's change set), following rhiza guide plan-reviser 'Applying the arbiter's rulings' literally: rhiza run status --json gave the serving-draft question change=<draft> and the other none; rhiza show <decision> --change <it> exit 0 for both (bare rhiza show exit 4, as the guide now avoids); rhiza apply set export-csv --change <draft> and keep-flat --record exit 0; decided-by links via rhiza apply --change <it> exit 0; rhiza run review-check --rulings passed 2 rulings applied; the run then went rulings visit 2 pass, gate, open-pr, merge, succeeded. Scoped tests: internal/api TestRunInfoAcceptsAndQuestions|TestRunReviewCheckRulings, internal/guide TestPlanReviserText, cmd/rhiza describe-golden, plan-workflow-autonomous-left, plan-workflow-autonomous-unapplied all ok.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
