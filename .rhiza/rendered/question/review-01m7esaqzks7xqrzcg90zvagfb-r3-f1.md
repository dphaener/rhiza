# review-01m7esaqzks7xqrzcg90zvagfb-r3-f1

question · `01M7ETG4ZKGAHYDAE4YRNCNADX`

## statement

Plan review round 3 left finding f1 (contradiction on budgets-out-of-make-test): The leaf takes 'everything CI runs' out of make check's help text, AGENTS.md and README.md because make check will no longer run the timing budgets that CI's budgets job runs. That contradicts the live, owner-accepted Decision t1-d10-makefile-targets ('`check` = everything CI runs.'). No node in the plan supersedes or amends it: budgets-opt-in supersedes only budgets-run-alone. The model would keep a Decision the code no longer follows. Fix: have budgets-opt-in (or a new Decision) supersede or restate t1-d10's check clause. The critic confirmed it: Verified: t1-d10-makefile-targets (01M5SWRRHWVZC2PG6R7EPZ82SJ) is a live owner-accepted Decision whose statement ends '`check` = everything CI runs.', and rhiza show lists no links_in, so nothing supersedes it. The plan's leaf budgets-out-of-make-test drops the budgets from make check while ci.yml:187-206's budgets job runs make test-budgets, and it removes 'everything CI runs' from the Makefile, AGENTS.md and README.md. budgets-opt-in supersedes only budgets-run-alone and restates nothing about check. The model would keep a Decision the code contradicts, which is a contradiction finding, so major stands. How should it be settled?

## Links

### Out

- blocks [budgets-out-of-make-test](../capability/rhiza/integrity/self-check/budgets-out-of-make-test.md)

### In

- answers [check-is-ci-minus-budgets](../decision/check-is-ci-minus-budgets.md)

## History

- `01M7EQ3AKCN5WTRED0VFN3JJD5` Read budgets are opt-in benchmarks: the timing budgets (TestReadBudget, TestReadBudgetApply, TestBudget in internal/index and internal/slice, TestCatchUpBudget, TestDriverPollIsCheap, TestGeneratedModel) leave make test for an opt-in make test-budgets that CI's budgets job runs, and the tests that wait out the index's 2 s racy window backdate mtimes with os.Chtimes instead of sleeping 2.2 s. Evidence: the budget phase is 188 s even warm.
