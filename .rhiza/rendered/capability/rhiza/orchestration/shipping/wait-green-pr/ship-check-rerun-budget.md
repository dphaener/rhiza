# ship-check-rerun-budget

capability · `01M6WJ6ZC52EWCYMSA2WW63P4V` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/api with a fake gh shows ship check with --reruns 2 rerunning a head's failed run twice, then failing on the third red, the ship record holding the count after each rerun; without the flag a head is rerun exactly once, as today; a new head resets the count; --reruns 0 is a usage error; the testscripts run-ship-rerun-pass and run-ship-rerun-fail still pass unchanged; describe-golden is regenerated for the new flag; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-rerun|describe-golden)' only

## statement

rhiza ship check takes --reruns <n> (n >= 1, default 1, today's once-per-head behaviour): when a settled check of a head fails, it reruns the failed jobs of that head's failed Actions runs up to n times per head, recording each rerun and the head's rerun count in the ship record before calling gh run rerun --failed (so a crash costs a missed rerun, never an extra one), and reports not settled while reruns remain; a head whose n reruns are spent fails as today, and --reruns 0 or a negative n is a usage error

## Links

### Out

- decided-by [plan-ci-flake-reruns](../../../../../decision/plan-ci-flake-reruns.md)

### In

- depends-on [plan-run-reruns-ci-flakes](../../../workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)
- verifies [ship-check-rerun-budget-verified](../../../../../evidence/ship-check-rerun-budget-verified.md)
- verifies [ship-check-rerun-budget-qa-1](../../../../../evidence/ship-check-rerun-budget-qa-1.md)
- verifies [pr-386-fix-the-plan-workflow-ends-a](../../../../../evidence/pr-386-fix-the-plan-workflow-ends-a.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
