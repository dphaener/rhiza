# failure-case-words

decision · `01M7KJ4SQDPV7ZY5V2KS6J432H`

## reasoning

A deterministic test needs a marker; the model already uses this one, and judging the case's substance stays with the reviewer

## statement

An acceptance shows its failure case with the words `failure case` (any case), as landed acceptances already write `(failure case)`; review-check tests only that the words or a failure_waiver are there, and the plan reviewer's no-failure-case kind judges whether the case is a real failure path

## Links

### Out

- derived-from [leaf-contract-plan](../source/leaf-contract-plan.md)

### In

- decided-by [leaf-contract-review-check](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-review-check.md)
- decided-by [planner-leaf-contract](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-leaf-contract.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
