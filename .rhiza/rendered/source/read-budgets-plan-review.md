# read-budgets-plan-review

source · `01M7ESG2MD03F9X7CK51TP8X7K`

## quotes

```json
[]
```

## raw_sha256

734148ec92c79c2adcf187d296860360f414e15c6e6afe59a5da131725bee231

## session

65d7da12-c81f-42cb-beaf-01a1dd95c54e

## summary

Autonomous plan run (owner away, orchestrator as arbiter). The orchestrator asked the planner to review and complete the two drafted leaves of milestone 38's read-budgets change set without widening scope, and to say how much mtime backdating cuts CI's budgets job (median 233 s, over the 3-minute target). Research: the internal/api budget tests wait out the index's 2 s racy window with 44 sleeps of 2.2 s (23 in TestReadBudgetApply, 21 in TestReadBudget's changeShowBudget and driverPollBudget), about 97 s of the 160 s the package's budgets take on the VM; CI's internal/api budgets step takes about 141 s. The index's racy check reads only modification times (internal/index/listing.go); the change time goes only into the fingerprint, so os.Chtimes backdating is sound and the drafted ctime contingency is settled. TestGeneratedModel left with internal/tui/modelapp and no longer exists. The drafted grep for 'time.Sleep(2' would also match time.Sleep(2 * pollDebounce), and budget_test.go and watchdog_pollbehind_test.go hold no racy-window sleep. Decided, assumed by the planner (would have escalated): a 75 s target for the internal/api budgets after backdating, an expected budgets-job median of about 135 s, and dropping TestGeneratedModel from BUDGET_TESTS.

## title

Read budgets opt-in plan review, autonomous

## Links

### In

- derived-from [backdate-cut-target](../decision/backdate-cut-target.md)
- derived-from [drop-generated-model-budget](../decision/drop-generated-model-budget.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
