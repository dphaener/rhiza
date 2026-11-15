# every-budget-opt-in

capability · `01M7FT14V71K8ZJGQ2P4NXF2FG`

## acceptance

go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^(TestMakeTestShardTargets|TestMakeTestTimeout)$' passes with the regex it pins naming both tests; make -n test-budgets runs them and make -n test, make -n test-api SHARD=1/2 and make -n test-rest skip them; every test that takes budgetLock is named by BUDGET_TESTS (grep -rl 'budgetLock(t)' --include=*_test.go . | xargs grep -B40 'budgetLock(t)' | grep -o 'func Test[A-Za-z]*' lists no name the regex misses); go test -list '.*' ./internal/api | grep -c '^Test' is unchanged. Touches Makefile, cmd/rhiza/makeupdate_test.go, and tools/gate only if it copies the regex.

## statement

Every timing budget is opt-in: BUDGET_TESTS also names TestInboxGateQuestionsBudget and TestRootIdlePollBudget, so make test and CI's api and rest jobs skip them and make test-budgets runs them with the others.

## Links

### Out

- decided-by [idle-and-gate-budgets-opt-in](../decision/idle-and-gate-budgets-opt-in.md)

### In

- verifies [every-budget-opt-in-verified](../evidence/every-budget-opt-in-verified.md)
- verifies [every-budget-opt-in-qa-1](../evidence/every-budget-opt-in-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
