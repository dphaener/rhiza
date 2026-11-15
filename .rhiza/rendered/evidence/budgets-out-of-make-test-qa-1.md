# budgets-out-of-make-test-qa-1

evidence · `01M7F0RZKVFNE4MEX2Q1S8S7PY` · actor [claude-code-8bb3496b](../actor/claude-code-8bb3496b.md)

## result

PASS

## statement

QA: make -n test shows both go tests with -skip of the five-name BUDGET_TESTS regex, no -run of it, no 'Timing budgets' phase or RHIZA_BUDGET_LOG; make -n check likewise (2 skip lines, 0 budget -run); make -n test-budgets runs -count=1 -p 1 -run BUDGET_TESTS over ./internal/api ./internal/index ./internal/slice and prints the timings; go test -list finds every name (TestReadBudget, TestReadBudgetApply, TestDriverPollIsCheap in api, TestBudget and TestCatchUpBudget in index, TestBudget in slice; TestGeneratedModel exists nowhere); test-scripts/test-api/test-rest dry-runs skip the budgets; SHARD=0/2 fails cleanly; ci.yml budgets job runs make test-budgets; make help, AGENTS.md and README.md say make check is what CI runs but make test-budgets; grep 'everything CI runs' Makefile AGENTS.md README.md prints nothing; TestMakeTestTimeout, TestMakeTestShardTargets, TestMakeUpdateDryRun pass; failure case: in a scratch worktree with a budget -run added to make test's recipe, TestMakeTestTimeout fails with 'make -n test runs the timing budgets'.

## Links

### Out

- verifies [budgets-out-of-make-test](../capability/rhiza/integrity/self-check/budgets-out-of-make-test.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
