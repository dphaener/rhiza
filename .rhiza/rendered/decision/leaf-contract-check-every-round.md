# leaf-contract-check-every-round

decision · `01M7KJ4SSK00BB6VTBGS78MXBB`

## reasoning

A reviser can always add a missing field or failure case, and a leaf without its contract must not reach a build

## statement

review-check's leaf-contract check fails at every review round, its problems never left as questions at the round cap

## Links

### Out

- derived-from [leaf-contract-plan](../source/leaf-contract-plan.md)

### In

- decided-by [leaf-contract-review-check](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-review-check.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
