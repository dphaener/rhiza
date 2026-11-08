# ship-check-rerun-budget-qa-1

evidence · `01M6XJKBHTZNKZRJ3B8YA9MX7M` · actor [claude-code-14b9e333](../actor/claude-code-14b9e333.md)

## result

PASS

## statement

QA cycle 1, by hand with the built binary (go build -trimpath ./cmd/rhiza) in a scratch repo with a bare origin, a change set branch cs/qa-cs and a fake gh whose PR check 'test (4/5)' (Actions run 42) stays FAILURE. rhiza ship check --reruns 2 called four times on one head: calls 1 and 2 printed 'not settled: reran the failed jobs of run 42' and exited 1, with the ship record's rerun_count at 1 then 2 (written before gh run rerun); calls 3 and 4 failed with no further rerun; gh.log held exactly two 'run rerun 42 --failed'. A new head with no flag was rerun exactly once (rerun_count reset to 1), then failed. Flipping the check to SUCCESS passed ('1 check passed', exit 0). --reruns 0, --reruns -1 and api call ship.check with {"reruns":0} each gave the usage error '--reruns takes a whole number of at least 1' (exit 2); --reruns abc gave 'flag --reruns takes an integer' (exit 2). Scoped tests passed: go test ./internal/api -run Ship, ./internal/workflow/..., and TestScripts (plan-workflow|run-ship-rerun|ship-check-ci-test-findings|describe-golden).

## Links

### Out

- verifies [ship-check-rerun-budget](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
