# budgets-opt-in

decision · `01M7EQCMV704AS2328FNE7RTMV`

## quote

those  are all fine. continue without me

## quote_where

owner, 2026-11-15T02:44:54Z, answering 'Budgets as opt-in: moving read budgets out of make test means CI runs them only in its budgets job, not on every local run. OK?' among the proposal's four questions

## statement

The timing budgets leave make test: make test skips them, and they run only in the opt-in make test-budgets, alone and one package at a time, which CI's budgets job runs on every pull request; the orchestrator's proposal, approved by the owner

## Links

### Out

- supersedes [budgets-run-alone](budgets-run-alone.md)
- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [budgets-measure-in-quiet](../capability/rhiza/views/fast-reads/budgets-measure-in-quiet.md)
- decided-by [make-test-shard-targets](../capability/rhiza/integrity/self-check/make-test-shard-targets.md)
- decided-by [budgets-out-of-make-test](../capability/rhiza/integrity/self-check/budgets-out-of-make-test.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
