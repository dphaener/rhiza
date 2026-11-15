# shard-targets-acceptance-follows-budget-move

decision · `01M7ETKP6TY9S36RMNR24HCTE2`

## reasoning

The change set's intent moves the timing budgets out of make test into make test-budgets, and leaf budgets-out-of-make-test requires TestMakeTestTimeout's expected commands to change (cmd/rhiza/makeupdate_test.go:81-90 asserts the budget -run line and the 'timing budgets' lock phase). The node's statement was already amended to say make test runs the suite but the budgets, so only the word 'unchanged' contradicts the intent, the leaf and the built code. The precedent in this plan is budgets-measure-in-quiet, whose acceptance was amended to point at the leaf; doing the same keeps the verified capability's acceptance true after the build without widening scope or dropping any existing check.

## statement

make-test-shard-targets' acceptance is amended in this change set, as budgets-measure-in-quiet's was: its last clause 'TestMakeTestTimeout still passes unchanged for make test' becomes 'TestMakeTestTimeout passes with its expected commands updated so make -n test shows the budget tests skipped and no Timing budgets phase (@budgets-out-of-make-test)'; the rest of the acceptance (the per-target make -n checks, test-api's shard partition and test-budgets printing the timings) stays as is.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7esaqzks7xqrzcg90zvagfb-r3-f2](../question/review-01m7esaqzks7xqrzcg90zvagfb-r3-f2.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
