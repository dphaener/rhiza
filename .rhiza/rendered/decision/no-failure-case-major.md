# no-failure-case-major

decision · `01M7KJ4SVS7PSKZMRDNAFXT04F`

## reasoning

100 of 160 verifier findings were acceptance points unmet or untested (prompt-gap rec 7): an untested failure path is an acceptance gap, the same severity as missing-acceptance

## statement

A no-failure-case plan-review finding is major, so review-check fails on a confirmed one before round 3

## Links

### Out

- derived-from [leaf-contract-plan](../source/leaf-contract-plan.md)

### In

- decided-by [no-failure-case-kind](../capability/rhiza/workflows/built-in-workflows/plan-workflow/no-failure-case-kind.md)
- decided-by [plan-reviewer-leaf-contract](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-leaf-contract.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
