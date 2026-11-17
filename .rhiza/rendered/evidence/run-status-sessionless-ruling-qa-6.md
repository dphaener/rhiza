# run-status-sessionless-ruling-qa-6

evidence · `01M7MW73BX0ZXQQC4B8AXSYCG2` · actor [claude-code-f91fbff2](../actor/claude-code-f91fbff2.md)

## result

PASS

## statement

QA cycle 6 at 85cf15dd8, scratch build go build -trimpath -o $SCRATCH/rhiza ./cmd/rhiza run by full path. In the kept work dir of plan-workflow-autonomous-unapplied, rhiza run status <run> printed both session-less review questions as "— decision settle-..." with no "not yet delivered"; --json questions had no session field, the serving-draft question carried change=<draft>, the landed-node one none. Scoped tests all ok: internal/api TestRunReviewCheckRulings|TestReviewCheck|TestRunInfoAcceptsAndQuestions (incl. TestReviewCheckQuestionInNodeChangeSet), internal/run TestCommandAsked|TestArbitrate (incl. TestCommandAskedArbitratedInChange), internal/guide TestPlanReviserText, internal/workflow TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview, cmd/rhiza scripts describe-golden, plan-workflow, plan-workflow-pr, plan-workflow-autonomous, -left, -unapplied, plan-review-visit-round, plan-milestone-backlog, run-review-check, -landed, -handles, run-arbitrate-asked.

## Links

### Out

- verifies [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
