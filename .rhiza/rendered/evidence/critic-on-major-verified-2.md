# critic-on-major-verified-2

evidence · `01M7BS33ARJC0YGE2MFVXY56DT` · actor [claude-code-d1aefca4](../actor/claude-code-d1aefca4.md)

## result

PASS

## statement

Verified attempt 2 at 5893bce13. Re-ran the Checks: block: go test -trimpath ./internal/api -run 'TestReviewCheckConductor$' fails (exit 1, 'no rounds: <nil>') with run_review.go reverted to HEAD~1 and passes at HEAD; go test ./internal/api -run 'Review|Conductor|Run|Findings|Verdict|Session' exit 0; the seven cmd/rhiza TestScripts (run-ship-needs-leaves, build-leaves-status, workflow-invoke, run-review-check, plan-workflow, plan-workflow-pr, describe-golden) exit 0; go test ./internal/guide/... ./internal/run/... ./internal/workflow/... exit 0; rhiza ship generate and rhiza render exit 0 with no tree diff; gofmt and go vet clean on internal/api. Acceptance go test -trimpath ./internal/api -run TestReviewCheckConductor passes: run journaled under conduct (plan, code, qa) and review-check run from the later gate command step; minor rounds pass, blocker with no verdict and confirmed major fail naming them, downgraded blocker and dropped major pass, only the latest round counts, no 'has no review round yet' error; and now a declared check with no round fails naming it. Diff touches only internal/api/run_review.go, run_review_test.go and the builder's actor record; no stubs or TODOs.

## Links

### Out

- verifies [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
