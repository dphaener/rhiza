# budget-tests-backdate-mtimes

capability · `01M7EQCMXHXRXDVGNJQEQBADVY` · parent [self-check](../self-check.md)

## acceptance

grep -n 'time.Sleep(2200' internal/api/*_test.go finds nothing (time.Sleep(2 * pollDebounce) stays); go test -trimpath -count=1 -p 1 -run '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' ./internal/api passes with the budget limits unchanged and finishes within 75 s (2026-11-15 on the VM: 160 s), measured with /usr/bin/time -f '%e s'; failure case: a test in internal/api/backdate_test.go shows index.LogStamp reports racy for a log folder just written, and not racy once the helper has backdated it past the window. Expected CI effect: the budgets job's internal/api step drops from about 141 s to about 45 s, taking the job's median from 233 s to about 135 s (@backdate-cut-target). Touches internal/api/write_budget_test.go, change_show_budget_test.go, run_model_cache_test.go and the new internal/api/backdate_test.go (pattern: internal/index/freshness_test.go). Must not change racyWindow, the index's racy check or a budget limit.

## statement

Tests that need the log's folders older than the index's 2 s racy window (internal/index/listing.go, whose racy check reads only modification times; the change time goes only into the fingerprint) backdate those folders' modification times with os.Chtimes through one shared internal/api test helper instead of sleeping 2.2 s, so TestReadBudgetApply and TestReadBudget's changeShowBudget and driverPollBudget measure and check the same things without their 44 sleeps (about 97 s).

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [backdate-cut-target](../../../../decision/backdate-cut-target.md)

### In

- verifies [budget-tests-backdate-mtimes-verified](../../../../evidence/budget-tests-backdate-mtimes-verified.md)
- verifies [budget-tests-backdate-mtimes-qa-1](../../../../evidence/budget-tests-backdate-mtimes-qa-1.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
