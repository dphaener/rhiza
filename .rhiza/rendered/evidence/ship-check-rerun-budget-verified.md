# ship-check-rerun-budget-verified

evidence · `01M6XFZFKEJEZ5974VSSQM4VBJ` · actor [claude-code-e65b59f5](../actor/claude-code-e65b59f5.md)

## result

PASS

## statement

Verifier at 83447746. Read the diff: ShipCheckRequest gains --reruns (ShipReruns; an explicit 0 decodes to -1 so 0 and negatives are usage errors, absent means 1); shipRerunFailed counts reruns per head (RerunCount, legacy states with Reruns count as 1), resets on a new head, and writes the state file with the new count before calling gh run rerun --failed; the --main dispatch path gets the same budget. go test -trimpath ./internal/api/ -run TestShip -count=1: ok (TestShipCheckRerunBudget: --reruns 2 reruns run 42 twice with RerunCount 1 then 2 in the ship record, the third red settles failed with no third gh rerun, a new head reruns again with count reset to 1; TestShipCheckRerunDefaultOnce: no flag reruns once, then fails; TestShipCheckRerunsUsage: reruns 0 and -1 give CodeUsage). A scratch build run as rhiza ship check --reruns 0 and --reruns -2 both exit 2 with a usage message, so the argv path goes through the same decoding. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-rerun|describe-golden)': run-ship-rerun-pass, run-ship-rerun-fail (both unchanged by the diff) and describe-golden PASS; the golden carries the new reruns flag. A full ./internal/api run with Budget tests skipped had 0 failing tests; the package was marked FAIL only by the post-run guard on ~/code/.rhiza-root/usage.live.json (the live root daemon's usage file, unrelated to this diff), and TestReadBudget/TestReadBudgetApply timing failures came from load 11-17 on 4 cores and do not touch ship code.

## Links

### Out

- verifies [ship-check-rerun-budget](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
