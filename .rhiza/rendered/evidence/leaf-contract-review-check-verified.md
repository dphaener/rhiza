# leaf-contract-review-check-verified

evidence · `01M7KRZM24ACRZ8E6WJ1NBH55Y` · actor [claude-code-b6c0168b](../actor/claude-code-b6c0168b.md)

## result

PASS

## statement

At f326fd5b5 (capped scope, -p 2): go test -trimpath ./internal/api exit 0; -run 'ReviewCheck' exit 0 (TestReviewCheckLeafContract covers a clean round passing, a leaf missing files failing with its handle and 'lacks files', no failure case failing and a failure_waiver passing, round 3 still failing with no asked question and round 3 the latest, a conductor run unaffected); go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-workflow|plan-review|autonomous-review|review-check)' exit 0; the other changed scripts describe-golden (no -update), plan-milestone-backlog and run-arbitrate-asked pass; docs/reference and the golden are current (no diff). Diff read: conductorReviewCheck untouched (dispatched before the new check); only ops of plan change sets are checked, so landed leaves the plan doesn't touch are skipped; the leaf-contract error returns before the round-cap question path; no TODOs. Caveat: each node is read in the view of the last plan change set touching it, so a parent created in one plan change set whose child is added in a later one could be flagged as a leaf.

## Links

### Out

- verifies [leaf-contract-review-check](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-review-check.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
