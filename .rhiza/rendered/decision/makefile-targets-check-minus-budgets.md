# makefile-targets-check-minus-budgets

decision · `01M7F13ZF4V6P9W83AV5MVABQB`

## reasoning

check-is-ci-minus-budgets commits this change set to restating t1-d10-makefile-targets with its targets unchanged plus test-budgets, and check as everything CI runs but the timing budgets, following the owner-approved budgets-opt-in Decision. No new owner words: the target list is t1-d10's, the exception is check-is-ci-minus-budgets'.

## statement

Makefile targets: build, install (built in place into `$(PREFIX)`, default `~/bin`, never cp), test, test-budgets, vet, fmt, fmt-check, check, clean and help. `check` = everything CI runs except the timing budgets, which run only in the opt-in `make test-budgets` (CI's budgets job).

## wording

(assumed, would have escalated) restates the owner-accepted t1-d10-makefile-targets with the check clause Decision check-is-ci-minus-budgets settles

## Links

### Out

- decided-by [check-is-ci-minus-budgets](check-is-ci-minus-budgets.md)
- supersedes [t1-d10-makefile-targets](t1-d10-makefile-targets.md) — restated per Decision check-is-ci-minus-budgets

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
