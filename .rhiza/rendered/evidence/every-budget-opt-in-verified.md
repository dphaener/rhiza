# every-budget-opt-in-verified

evidence · `01M7FWCFSXEMWWHCW62BW8PV2A` · actor [claude-code-7633a7b5](../actor/claude-code-7633a7b5.md)

## result

PASS

## statement

At 4bf1a438b: go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^(TestMakeTestShardTargets|TestMakeTestTimeout)$' passes (both PASS), both pinning the regex that names TestInboxGateQuestionsBudget and TestRootIdlePollBudget. make -n test-budgets runs -run with the 7-name regex over BUDGET_PKGS (both new tests live in internal/api); make -n test and test-rest -skip it, make -n test-api SHARD=1/2 drops them via grep -Ev. grep -rl 'budgetLock(t)' ... | grep -o 'func Test[A-Za-z]*' lists TestBudget, TestCatchUpBudget, TestInboxGateQuestionsBudget, TestReadBudget, TestReadBudgetApply, TestRootIdlePollBudget: all in BUDGET_TESTS. go test -list '.*' ./internal/api | grep -c '^Test' = 762, internal/ and tools/ untouched (tools/gate does not copy the regex). Diff touches only Makefile, cmd/rhiza/makeupdate_test.go and .rhiza records; no TODOs.

## Links

### Out

- verifies [every-budget-opt-in](../capability/every-budget-opt-in.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
