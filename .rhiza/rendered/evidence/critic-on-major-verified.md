# critic-on-major-verified

evidence · `01M7BJ985VAHWGVZP4TC4RM8QS` · actor [claude-code-2c122a02](../actor/claude-code-2c122a02.md)

## result

PASS

## statement

At d483f08f0, re-ran the Checks: block: go test -trimpath ./internal/api -run TestReviewCheckConductor (exit 0, PASS), go test -trimpath ./internal/api -run 'TestReviewCheck|TestRun|TestPlanReview' (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked diff). Read the diff: run.review-check now finds the workflow step declaring checks (current step first, else any step) via checksDeclaringStep and, in conductorReviewCheck, counts each declared check's latest round under that step through Round.Standing (blocker/major with no verdict or confirmed stand; downgraded below major or dropped do not), failing with an error listing check and finding, never taking the plan-review path. TestReviewCheckConductor runs from the gate command step after a conduct step declaring plan, code, qa and covers minor/nit pass, blocker with no verdict fails naming code and the finding, confirmed major fails, downgraded blocker and dropped major pass, only the latest round counts (plan major then clean round), and asserts no 'has no review round yet' error. Diff touches only internal/api/run_review.go, run_review_test.go plus the builder's actor record and its rendering; no stubs or TODOs; fakes only.

## Links

### Out

- verifies [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
