# pr-386-fix-the-plan-workflow-ends-a

evidence · `01M6XSQJ73P580E9TTBGXEBZG5` · actor [dphaener](../actor/dphaener.md)

## head

f1e55734028892392d845d02f331f53a39ab0a7f

## pr

https://github.com/dphaener/rhiza/pull/386

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #386 (https://github.com/dphaener/rhiza/pull/386) passed CI at f1e5573

## status

verified

## Links

### Out

- verifies [ship-check-rerun-budget](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)
- verifies [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)
- verifies [plan-run-reruns-ci-flakes](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
