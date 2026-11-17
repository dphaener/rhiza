# review-question-in-node-change-set

capability · `01M7MR5TM3RWMPVJ2AND27RQH4`

## acceptance

go test -trimpath ./internal/api -run 'TestReviewCheckQuestionInNodeChangeSet' passes, covering, with the run's change set open: a confirmed round-3 finding on a leaf a serving draft created puts its question (and blocks link) in that draft and the asked entry's change names the draft; a finding on a node the run's change set created stays in the run's change set; and, with the run's change set landed, a finding on a node no in-flight plan change set created is still a record write. go test -trimpath ./internal/run -run 'TestCommandAskedArbitratedInChange' passes: the arbiter for a session-less asked entry naming change set X starts with RHIZA_CHANGE=X, and one naming none starts with the run's change set. go test -trimpath ./internal/run -run 'TestCommandAsked|TestArbitrate' and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-review-check|run-review-check-landed|run-review-check-handles|run-arbitrate-asked)$' still pass.

## failure_waiver

Placement only: the failure paths are unchanged and already covered (an arbiter that cannot decide fails the run naming the question, TestArbitrateFailsWithoutAnswer); a ruling left unapplied failing the run is plan-workflow-applies-rulings' failure case.

## files

internal/api/run_review.go, internal/api/run_review_test.go, internal/run/engine.go, internal/run/engine_test.go, internal/run/journal.go

## must_not

Change where questions on landed nodes or the run's own nodes go; change what review-check --rulings reads; change how a question a session asked is arbitrated or delivered; let the arbiter write anything but its Decision.

## statement

rhiza run review-check puts a round-3 question on a finding whose node an in-flight plan change set created (a serving draft's leaf) into that change set, whether or not the run's change set has landed (reviewQuestionTarget no longer returns the run's change set first while it is open), and journals the change set it wrote the question into on the asked entry's omitempty change field. The engine starts the arbiter for a session-less asked entry with RHIZA_CHANGE set to the entry's change, falling back to the run's change set when the entry names none, so the question, the arbiter's Decision and the leaf share one change set's view. Findings on landed nodes or on nodes of the run's own change set keep today's placement.

## Links

### Out

- decided-by [review-ruling-in-node-change-set](../decision/review-ruling-in-node-change-set.md)
- decided-by [plan-rulings-applied-before-gate](../decision/plan-rulings-applied-before-gate.md)

### In

- depends-on [plan-workflow-applies-rulings](rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)
- verifies [review-question-in-node-change-set-verified-2](../evidence/review-question-in-node-change-set-verified-2.md)
- verifies [review-question-in-node-change-set-qa-1](../evidence/review-question-in-node-change-set-qa-1.md)
- verifies [review-question-in-node-change-set-qa-2](../evidence/review-question-in-node-change-set-qa-2.md)
- verifies [review-question-in-node-change-set-qa-3](../evidence/review-question-in-node-change-set-qa-3.md)
- verifies [review-question-in-node-change-set-qa-4](../evidence/review-question-in-node-change-set-qa-4.md)
- verifies [review-question-in-node-change-set-qa-5](../evidence/review-question-in-node-change-set-qa-5.md)
- verifies [review-question-in-node-change-set-qa-6](../evidence/review-question-in-node-change-set-qa-6.md)
- verifies [review-question-in-node-change-set-qa-7](../evidence/review-question-in-node-change-set-qa-7.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
