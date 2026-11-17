# leaf-contract-review-check

capability · `01M7KJ4T26ZW4X7KSWP850DHDF` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/api -run 'ReviewCheck' passes with tests that a clean round over leaves with all fields and a failure case passes; (failure case) a leaf missing files fails with an error naming the leaf and files; (failure case) a leaf whose acceptance has no failure case fails, and the same leaf with a failure_waiver passes; (failure case) the leaf-contract failure still fails at round 3 with no question added; and a conductor-run review-check is unaffected; go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-workflow|plan-review|autonomous-review|review-check)' passes.

## statement

`rhiza run review-check` on the plan-review path (no checks step; not --fresh or --findings) fails, at every round, with an error listing each live leaf (a capability with no live child capability) that a plan change set creates or changes and that lacks a statement, acceptance, files or must_not, or whose acceptance holds no `failure case` (any case) and that has no failure_waiver, naming the leaf and what it lacks; it then checks the round as today, and passes only when both pass. The conductor path is unchanged. Files: internal/api/run_review.go (or a new internal/api/run_review_leaves.go beside it), internal/api/run_review_test.go, and the plan-review fixtures in internal/api tests and cmd/rhiza/testdata/script (plan-workflow*.txtar, plan-review-*.txtar, autonomous-review.txtar and others that run review-check) given the fields so they pass. Must not: change conductorReviewCheck, turn a leaf-contract problem into a question at the round cap, or check landed leaves the plan does not touch.

## Links

### Out

- depends-on [leaf-contract-fields](leaf-contract-fields.md)
- decided-by [failure-case-words](../../../../../decision/failure-case-words.md)
- decided-by [leaf-contract-check-every-round](../../../../../decision/leaf-contract-check-every-round.md)

### In

- depends-on [plan-reviewer-leaf-contract](plan-reviewer-leaf-contract.md)
- verifies [leaf-contract-review-check-verified](../../../../../evidence/leaf-contract-review-check-verified.md)
- verifies [leaf-contract-review-check-qa-1](../../../../../evidence/leaf-contract-review-check-qa-1.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
