# review-01m7esaqzks7xqrzcg90zvagfb-r3-f2

question · `01M7ETG50SJ5FHK62K1ETXAZPK`

## statement

Plan review round 3 left finding f2 (contradiction on make-test-shard-targets): The plan updates this verified capability's statement but keeps its acceptance, which still says 'TestMakeTestTimeout still passes unchanged for make test'. Leaf budgets-out-of-make-test requires TestMakeTestTimeout's expected commands to change: cmd/rhiza/makeupdate_test.go:83-87 asserts the 'Timing budgets' phase and its -run line, which leave make test. Once the leaf is built, this node's acceptance contradicts the leaf and the code. Its acceptance needs updating as budgets-measure-in-quiet's was. The critic confirmed it: Verified: make-test-shard-targets (verified) has its statement amended by this plan, but its acceptance is unchanged and still says 'TestMakeTestTimeout still passes unchanged for make test'. cmd/rhiza/makeupdate_test.go:81-90 asserts the budgets -run line and the lock "the timing budgets" phase in make -n test, and budgets-out-of-make-test's acceptance requires those expected commands updated and no 'Timing budgets' phase in make test. Once the leaf is built, this node's acceptance contradicts the leaf and the code, as budgets-measure-in-quiet's did before it was amended. Major stands. How should it be settled?

## Links

### Out

- blocks [make-test-shard-targets](../capability/rhiza/integrity/self-check/make-test-shard-targets.md)

### In

- answers [shard-targets-acceptance-follows-budget-move](../decision/shard-targets-acceptance-follows-budget-move.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
