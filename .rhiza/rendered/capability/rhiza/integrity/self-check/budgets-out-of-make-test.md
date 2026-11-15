# budgets-out-of-make-test

capability · `01M7EQDFY7GKPF549RAPQE285W` · parent [self-check](../self-check.md)

## acceptance

make -n test shows go test with -skip of the budget tests, no -run of them and no 'Timing budgets' phase; make -n test-budgets shows them run -count=1 -p 1 over ./internal/api ./internal/index ./internal/slice; BUDGET_TESTS is ^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap)$ and each name matches a test go test -list finds in those packages (@drop-generated-model-budget); go test -trimpath ./cmd/rhiza -run '^(TestMakeTestTimeout|TestMakeTestShardTargets|TestMakeUpdateDryRun)$' passes with its expected commands updated; .github/workflows/ci.yml's budgets job still runs make test-budgets; make help shows make check without 'everything CI runs' and naming make test-budgets as the CI run it leaves out, and AGENTS.md's Commands block line for make check and README.md's Development section say the same (no 'everything CI runs' left in any of them: grep -n 'everything CI runs' Makefile AGENTS.md README.md prints nothing); failure case: TestMakeTestTimeout or TestMakeTestShardTargets fails when make test's recipe runs a budget test. Speedup measure (orchestrator, quiet VM): a warm /usr/bin/time -f '%e s' make test, against 254 s, expected about 70 s since the budget phase took 188 s, with no 'Timing budgets' phase in its output. Touches Makefile (test's and check's help text and check's comment), cmd/rhiza/makeupdate_test.go, AGENTS.md (the '(and the timing budgets)' lock sentence, make test's description and the Commands block's 'make check # everything CI runs' line), README.md (the Development section's '`make check` runs everything CI runs.' sentence). Must not drop a budget test from test-budgets or change a budget limit.

## statement

make test and make check no longer run the timing budgets: the budget phase leaves make test's recipe, the budgets run only in the opt-in make test-budgets (alone, -count=1 -p 1, printing the read-budget timings), which CI's budgets job runs on every pull request, and AGENTS.md, README.md and the Makefile's comments say so.

## Links

### Out

- decided-by [budgets-opt-in](../../../../decision/budgets-opt-in.md)
- decided-by [drop-generated-model-budget](../../../../decision/drop-generated-model-budget.md)

### In

- blocks [review-01m7esaqzks7xqrzcg90zvagfb-r3-f1](../../../../question/review-01m7esaqzks7xqrzcg90zvagfb-r3-f1.md)
- verifies [budgets-out-of-make-test-verified](../../../../evidence/budgets-out-of-make-test-verified.md)
- verifies [budgets-out-of-make-test-qa-1](../../../../evidence/budgets-out-of-make-test-qa-1.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
