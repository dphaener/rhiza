# backdate-cut-target

decision · `01M7ESH0B3D7D7V7RNSQFJM1S4`

## question

How much should the backdating leaf be expected to cut, locally and in CI's budgets job?

## reasoning

The drafted 60 s was below what removing every racy-window sleep can reach: 160 s minus about 97 s of sleeps leaves about 63 s of real work on the VM, so 75 s gives margin for load without allowing a sleep to stay. CI's internal/api budgets step is about 141 s; the same 97 s of wall-clock sleeps come off it, leaving about 45 s, so the job's 233 s median falls to about 135 s.

## statement

After backdating, internal/api's three budget tests run -count=1 -p 1 must finish within 75 s on the VM (from 160 s, of which about 97 s are 44 racy-window sleeps), and CI's budgets job is expected to drop from a 233 s median to about 135 s, under the 3-minute CI target

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [read-budgets-plan-review](../source/read-budgets-plan-review.md)

### In

- decided-by [budget-tests-backdate-mtimes](../capability/rhiza/integrity/self-check/budget-tests-backdate-mtimes.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
