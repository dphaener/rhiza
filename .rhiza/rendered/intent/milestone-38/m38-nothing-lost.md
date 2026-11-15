# m38-nothing-lost

intent · `01M7EQ752DW9BDSP53JN8H9BJC` · parent [milestone-38](../milestone-38.md)

## statement

No test is deleted or skipped to get there: the number of tests go test -list '.*' ./... lists across packages, and the number of cmd/rhiza/testdata/script/*.txtar files (396), do not drop below their count at commit 88516d956, except tests moved to opt-in benchmarks, which still run in CI's budgets job. Measure: go test -list '.*' ./... | grep -c '^Test' and ls cmd/rhiza/testdata/script/*.txtar | wc -l, on 88516d956 and on main.

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../decision/fast-suite-proposal-approved.md)
- decided-by [fast-suite-measuring-commands](../../decision/fast-suite-measuring-commands.md)

### In

- serves change set `01M7EQ3AKCN5WTRED0VFN3JJD5` (Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
