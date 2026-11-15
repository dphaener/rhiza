# check-is-ci-minus-budgets

decision · `01M7ETHS5ACFFFZWV0CH3CGX7M`

## reasoning

The change set's intent and the owner-approved budgets-opt-in Decision ('those are all fine. continue without me', answering 'moving read budgets out of make test means CI runs them only in its budgets job, not on every local run. OK?') already chose that make test, and so make check, skip the budgets while CI's budgets job runs them. The owner knowingly accepted the gap, so putting the budgets back into make check would undo the change set's 188 s saving. t1-d10 dates from milestone 1 and its check clause just has to follow the newer owner choice. The finding's own fix (supersede or restate t1-d10) is the right mechanism: a restated Decision keeps t1-d10's target list and records the one exception, and the leaf's doc edits (dropping 'everything CI runs', naming make test-budgets as what check leaves out) then match the model.

## statement

Settle it by amending t1-d10's check clause, not by keeping the budgets in make check: make check runs everything CI runs except the timing budgets, which run only in the opt-in make test-budgets (CI's budgets job). The plan adds a Decision restating t1-d10-makefile-targets with that check clause (targets unchanged, plus test-budgets) and links it superseding t1-d10-makefile-targets, so no live Decision says check = everything CI runs; budgets-out-of-make-test keeps its acceptance as written.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7esaqzks7xqrzcg90zvagfb-r3-f1](../question/review-01m7esaqzks7xqrzcg90zvagfb-r3-f1.md)

### In

- decided-by [makefile-targets-check-minus-budgets](makefile-targets-check-minus-budgets.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
