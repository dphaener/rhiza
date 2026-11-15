# budgets-measure-in-quiet

capability · `01M6HR4C6TN1MN8FF6C71MZ3F2` · parent [fast-reads](../fast-reads.md)

## acceptance

make -n test-budgets shows the timing budgets run alone, -count=1 -p 1, one package at a time, and make -n test runs none of them (@budgets-out-of-make-test)

## statement

The timing budgets are opt-in benchmarks: make test runs the suite with them skipped, and make test-budgets runs them alone, -count=1 -p 1 one package at a time, printing the read-budget timings; CI's budgets job runs make test-budgets on every pull request, so no budget test is lost.

## Links

### Out

- decided-by [budgets-run-alone](../../../../decision/budgets-run-alone.md)
- decided-by [budgets-opt-in](../../../../decision/budgets-opt-in.md)

### In

- verifies [pr-301-delete-the-old-rex-apps](../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6HR3VBT4DPCNS03S70QCAZF` Timing budgets measure in quiet: make check runs TestReadBudget, TestBudget and TestGeneratedModel on their own after the rest of the suite, one package at a time, so CI and the VM stop failing them on their neighbours' load
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
