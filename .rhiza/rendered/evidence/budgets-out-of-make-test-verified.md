# budgets-out-of-make-test-verified

evidence · `01M7EZ9DX9N6H244KS0V49E6BQ` · actor [claude-code-ca8456f0](../actor/claude-code-ca8456f0.md)

## result

PASS

## statement

At 857f53e67, re-ran the Checks block: go test -trimpath -p 2 ./cmd/rhiza -run '^(TestMakeTestTimeout|TestMakeTestShardTargets|TestMakeUpdateDryRun)$' exit 0 (all three PASS); failure case in a scratch worktree with the budget phase (Timing budgets echo, lock "the timing budgets", go test -run BUDGET_TESTS) restored to make test's recipe: TestMakeTestTimeout FAILs with 'make -n test runs the timing budgets', exit 1; make -n test exit 0, shows go test with -skip of the budget regex for both phases, no -run of it, no Timing budgets phase; make -n test-budgets exit 0, runs -count=1 -p 1 -run '^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap)$' ./internal/api ./internal/index ./internal/slice under the lock and prints the timings; go test -list of that regex over the three packages finds every name (TestGeneratedModel dropped per drop-generated-model-budget; no func TestGeneratedModel exists in the repo); make help exit 0, check reads 'Run what CI runs but make test-budgets (the timing budgets)'; grep -n 'everything CI runs' Makefile AGENTS.md README.md exit 1; grep make test-budgets .github/workflows/ci.yml exit 0 (line 206); rhiza ship generate exit 0 (0 changed); rhiza render exit 0. Diff read: only Makefile, AGENTS.md, README.md, cmd/rhiza/makeupdate_test.go and the builder's .rhiza log/actor files change; no budget limit touched, no TODOs. Speedup measure skipped per brief (orchestrator's).

## Links

### Out

- verifies [budgets-out-of-make-test](../capability/rhiza/integrity/self-check/budgets-out-of-make-test.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
