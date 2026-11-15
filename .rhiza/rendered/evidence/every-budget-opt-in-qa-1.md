# every-budget-opt-in-qa-1

evidence · `01M7G3RG7RSNC42P3SR9HP9BQ3` · actor [claude-code-82f1640c](../actor/claude-code-82f1640c.md)

## result

PASS

## statement

QA: make -n test-budgets runs BUDGET_TESTS (now incl. TestInboxGateQuestionsBudget, TestRootIdlePollBudget); make -n test and test-rest skip them; every budgetLock test is named; TestMakeTestShardTargets and TestMakeTestTimeout pass

## Links

### Out

- verifies [every-budget-opt-in](../capability/every-budget-opt-in.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
