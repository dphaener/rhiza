# m38-warm

intent · `01M7EQ750SRBDSEFNF18FJ9HQT` · parent [milestone-38](../milestone-38.md)

## statement

A warm make test (run again right after a passing one, with no change) finishes within 60 s of wall time on the 4-core VM (2026-11-15: 254 s). Measure: make test, then /usr/bin/time -f '%e s' make test.

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../decision/fast-suite-proposal-approved.md)
- decided-by [fast-suite-measuring-commands](../../decision/fast-suite-measuring-commands.md)

### In

- serves change set `01M7EQ3AKCN5WTRED0VFN3JJD5` (Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.)
- serves change set `01M7EQ42Y1M8ENZCQ202D761JM` (Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
