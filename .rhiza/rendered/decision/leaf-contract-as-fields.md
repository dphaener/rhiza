# leaf-contract-as-fields

decision · `01M7KJ4SNGBQ8RQ0V96MGSWZVM`

## reasoning

prompt-gap section 5 names Go model fields files and must-not; fields let review-check test them mechanically and a later parallel-leaf scheduler compare files, which acceptance prose would not

## statement

A leaf's contract is three capability fields, files (the paths it touches, line ranges allowed), must_not (what it must not do) and failure_waiver (why its acceptance has no failure case), not text inside acceptance

## Links

### Out

- derived-from [leaf-contract-plan](../source/leaf-contract-plan.md)

### In

- decided-by [leaf-contract-fields](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-fields.md)
- decided-by [planner-leaf-contract](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-leaf-contract.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
