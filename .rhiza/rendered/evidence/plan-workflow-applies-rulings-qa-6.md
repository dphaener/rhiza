# plan-workflow-applies-rulings-qa-6

evidence · `01M7MW73QB6V29359T1DMP3WTF` · actor [claude-code-f91fbff2](../actor/claude-code-f91fbff2.md)

## result

PASS

## statement

QA cycle 6 at 85cf15dd8, scratch build go build -trimpath -o $SCRATCH/rhiza ./cmd/rhiza run by full path. Scratch autonomous plan run went check (visit 3) pass -> rulings (visit 1) fail -> apply-rulings (visit 1) -> rulings (visit 2) fail -> run failed, step apply-rulings went over max_visits 1, no gate, draft and record change sets still draft. Kept plan-workflow-autonomous-unapplied work dir: same shape, run status names both unapplied Decisions. Scoped tests all ok: internal/api TestRunReviewCheckRulings|TestReviewCheck|TestRunInfoAcceptsAndQuestions (incl. TestReviewCheckQuestionInNodeChangeSet), internal/run TestCommandAsked|TestArbitrate (incl. TestCommandAskedArbitratedInChange), internal/guide TestPlanReviserText, internal/workflow TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview, cmd/rhiza scripts describe-golden, plan-workflow, plan-workflow-pr, plan-workflow-autonomous, -left, -unapplied, plan-review-visit-round, plan-milestone-backlog, run-review-check, -landed, -handles, run-arbitrate-asked.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
