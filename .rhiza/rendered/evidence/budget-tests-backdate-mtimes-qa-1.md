# budget-tests-backdate-mtimes-qa-1

evidence · `01M7F0RZRA7D696VSCHH42TWG6` · actor [claude-code-8bb3496b](../actor/claude-code-8bb3496b.md)

## result

PASS

## statement

QA: grep 'time.Sleep(2200' internal/api/*_test.go finds nothing, time.Sleep(2 * pollDebounce) stays (run_model_cache_test.go:220); no internal/index change, racyWindow still 2s, no budget limit in the diff; behind the machine lock (load 3.4) /usr/bin/time go test -trimpath -count=1 -p 1 -run '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap|TestBackdateLogPassesRacyWindow)$' ./internal/api passes in 67.74 s (under 75 s; TestReadBudget 60.2 s, Apply 3.5 s, DriverPoll 1.7 s); TestBackdateLogPassesRacyWindow passes: LogStamp racy right after a write, not racy after backdateLog.

## Links

### Out

- verifies [budget-tests-backdate-mtimes](../capability/rhiza/integrity/self-check/budget-tests-backdate-mtimes.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
