# budget-tests-backdate-mtimes-verified

evidence · `01M7EZVJBP2DP3DQ8QPPX3HFN1` · actor [claude-code-5132081b](../actor/claude-code-5132081b.md)

## result

PASS

## statement

At d284ba557 behind the machine lock with RHIZA_/SHIP_/REX_ unset: grep -n 'time.Sleep(2200' internal/api/*_test.go exit 1 (no match; time.Sleep(2 * pollDebounce) stays at run_model_cache_test.go:220); go vet ./internal/api exit 0; go test -trimpath -p 2 ./internal/api -run '^TestBackdateLogPassesRacyWindow$' -count=1 ok (failure case: index.LogStamp racy for a folder just written, not racy after backdateLog); /usr/bin/time -f '%e s' go test -trimpath -count=1 -p 1 -run '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' ./internal/api ok, test 73.19 s, wall 75.00 s (at the 75 s target, under load; was 160 s). Diff touches only the four named test files plus the builder's actor record; internal/index (racyWindow, racy check) and every budget limit unchanged. Measurement equivalence: takeStatsBy's racy check reads only log/ folder ModTimes against now-2s; backdateLog moves back exactly the folders with mtime within 2.2 s (the ones a 2.2 s sleep would age past the window) and leaves older folders' stats, hence any stored fingerprint, untouched; no other read-path cache compares log times to the clock.

## Links

### Out

- verifies [budget-tests-backdate-mtimes](../capability/rhiza/integrity/self-check/budget-tests-backdate-mtimes.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
