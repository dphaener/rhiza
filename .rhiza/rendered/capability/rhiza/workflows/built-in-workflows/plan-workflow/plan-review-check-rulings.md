# plan-review-check-rulings

capability · `01M7MPK6WHYWK078TGV14QZA33` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/api -run 'TestRunReviewCheckRulings' passes, covering: no asked entry passes; a passed ruling whose Decision a live node links decided-by passes (one case with a capability, one with a Decision node in a landed record); an open (unarbitrated) question passes; and go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden' passes with the golden regenerated for the new flag. Failure case: a passed ruling whose Decision has only its answers link fails with an error naming the question handle and the Decision handle, and the same with the only decided-by link retired.

## files

internal/api/run_review.go, internal/api/run_review_test.go, cmd/rhiza/testdata/script/describe-golden.txtar, docs/reference/cli.md

## must_not

Change what review-check does without --rulings, its conductor path (checksDeclaringStep) or the questions it adds at round 3; write to the model or the journal; count a question a session asked.

## statement

rhiza run review-check --rulings passes when every review question of the run (a session-less asked entry) whose arbitration passed has its answering Decision applied, a live node linking decided-by to it, and when there is no such question; otherwise it exits non-zero listing each unapplied ruling (question handle, Decision handle, the node of the finding it came from) with next rhiza run status. It judges and writes nothing.

## Links

### Out

- decided-by [plan-rulings-applied-before-gate](../../../../../decision/plan-rulings-applied-before-gate.md)
- decided-by [ruling-applied-means-decided-by](../../../../../decision/ruling-applied-means-decided-by.md)

### In

- depends-on [plan-workflow-applies-rulings](plan-workflow-applies-rulings.md)
- verifies [plan-review-check-rulings-verified](../../../../../evidence/plan-review-check-rulings-verified.md)
- verifies [plan-review-check-rulings-qa-1](../../../../../evidence/plan-review-check-rulings-qa-1.md)
- verifies [plan-review-check-rulings-qa-2](../../../../../evidence/plan-review-check-rulings-qa-2.md)
- verifies [plan-review-check-rulings-qa-3](../../../../../evidence/plan-review-check-rulings-qa-3.md)
- verifies [plan-review-check-rulings-qa-4](../../../../../evidence/plan-review-check-rulings-qa-4.md)
- verifies [plan-review-check-rulings-qa-5](../../../../../evidence/plan-review-check-rulings-qa-5.md)
- verifies [plan-review-check-rulings-qa-6](../../../../../evidence/plan-review-check-rulings-qa-6.md)
- verifies [plan-review-check-rulings-qa-7](../../../../../evidence/plan-review-check-rulings-qa-7.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
