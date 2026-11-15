# drop-generated-model-budget

decision · `01M7ESH0D6NJ2X7HZSR1XW9G58`

## question

BUDGET_TESTS names TestGeneratedModel, which no package holds any more: drop it or keep it?

## reasoning

A name in BUDGET_TESTS that matches no test keeps a dead entry in every -run and -skip; no budget is lost by dropping it, since the test it named no longer exists anywhere in the repo.

## statement

TestGeneratedModel is dropped from BUDGET_TESTS and from the tests that spell the budget regex: it lived in internal/tui/modelapp, which was deleted (PR #301), and no package holds it now

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [read-budgets-plan-review](../source/read-budgets-plan-review.md)

### In

- decided-by [budgets-out-of-make-test](../capability/rhiza/integrity/self-check/budgets-out-of-make-test.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
