# leaf-contract-review-check-qa-1

evidence · `01M7KW24XRKHA0QADGHVEVSXHA` · actor [claude-code-ba01156e](../actor/claude-code-ba01156e.md)

## result

PASS

## statement

QA with the built binary on a hold-workflow plan run: review-check failed exit 1 naming cap-b lacks files, must_not, a failure case; a parent with a live child was skipped and an uppercase FAILURE CASE passed; a waiver satisfied the failure case; whitespace-only must_not reported as lacking; round 2 printed both the leaf-contract and confirmed-finding failures; round 3 still failed with no question or asked entry; once fixed, round 4 passed leaving the confirmed finding as a question; retiring the child made the parent a leaf and review-check failed naming it.

## Links

### Out

- verifies [leaf-contract-review-check](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-review-check.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
