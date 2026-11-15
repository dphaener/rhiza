# idle-and-gate-budgets-opt-in

decision · `01M7FT14RK9GZJAT8R9QXZGQ9W`

## reasoning

Both measure time against a fixed limit and take the budget lock, like the five budgets #590 moved out of make test under fast-suite-proposal-approved, but the skip regex does not name them, so they run serially in make test and in CI's api shards (13.7 s and 9.6 s). They still run in CI's budgets job, so m38-nothing-lost holds (its exception for tests moved to opt-in benchmarks). The budgets job gains about 23 s, which the fresh CI build cache more than repays.

## statement

TestInboxGateQuestionsBudget and TestRootIdlePollBudget join the opt-in timing budgets (BUDGET_TESTS): make test, test-api and test-rest skip them and make test-budgets runs them

## wording

(assumed, would have escalated)

## Links

### Out

- relates-to [fast-suite-proposal-approved](fast-suite-proposal-approved.md) — Applies its opt-in budgets ruling to two budgets #590 missed
- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [every-budget-opt-in](../capability/every-budget-opt-in.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
