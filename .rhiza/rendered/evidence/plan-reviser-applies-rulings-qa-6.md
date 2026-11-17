# plan-reviser-applies-rulings-qa-6

evidence · `01M7MW73GZ62DSEJCN947A1DX4` · actor [claude-code-f91fbff2](../actor/claude-code-f91fbff2.md)

## result

PASS

## statement

QA cycle 6 at 85cf15dd8, scratch build go build -trimpath -o $SCRATCH/rhiza ./cmd/rhiza run by full path. Followed the guide section literally in a scratch reviser: for an oversized-leaf split in a serving draft, rhiza show <decision> --change <question change> exit 0, retire + two new leaves in the draft, decided-by links from each in the question change exit 0, --rulings then counted that ruling applied. The link-finding case could not be applied (see the review-question-in-node-change-set finding). Scoped tests all ok: internal/api TestRunReviewCheckRulings|TestReviewCheck|TestRunInfoAcceptsAndQuestions (incl. TestReviewCheckQuestionInNodeChangeSet), internal/run TestCommandAsked|TestArbitrate (incl. TestCommandAskedArbitratedInChange), internal/guide TestPlanReviserText, internal/workflow TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview, cmd/rhiza scripts describe-golden, plan-workflow, plan-workflow-pr, plan-workflow-autonomous, -left, -unapplied, plan-review-visit-round, plan-milestone-backlog, run-review-check, -landed, -handles, run-arbitrate-asked.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
