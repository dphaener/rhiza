# plan-review-check-rulings-verified

evidence · `01M7MQP26H683PE8YA6QQP5YNP` · actor [claude-code-91fd8633](../actor/claude-code-91fd8633.md)

## result

PASS

## statement

At 5752832ae: go test -p 2 -trimpath -count=1 ./internal/api -run 'TestRunReviewCheckRulings|TestReviewCheck' exit 0 (TestRunReviewCheckRulings plus the 6 existing TestReviewCheck* pass); go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden exit 0, and with -update leaves no diff; rhiza docs, rhiza ship generate, rhiza render (scratch build) exit 0, render touching only other sessions' uncommitted actors. The test covers no asked entry, open review questions, a session's question ignored, a capability linked decided-by, a landed record Decision linked decided-by, and both failure cases (answers link only; the only decided-by link retired) asserting the question handle, Decision handle and next rhiza run status R40. Diff read: --rulings returns before checksDeclaringStep, so review-check without it is unchanged; reviewRulingsCheck only reads (reviewQuestions filters Session==""); no stubs or TODOs; files match the leaf plus the builder's actor record and render.

## Links

### Out

- verifies [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
