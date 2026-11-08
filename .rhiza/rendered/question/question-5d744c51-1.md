# question-5d744c51-1

question · `01M6WJ10YJKJ6GG53EVNNFCYAX`

## session

5d744c51-b515-4500-b578-01a1b91df18d

## statement

Change set 01M6WHT (plan workflow ends a run on a CI flake) serves no intent yet. Which intent should it serve? / ship check already reruns a head's failed Actions jobs once (ci-rerun-once-per-head). In run 01M6VY5APJP9 that rerun fired and both flaky specs failed again; a third attempt passed. How should the plan workflow get its extra rerun? / When a check stays red after the reruns, how should ship check say the failure is unrelated to the plan? / Should a .rhiza-only PR whose app tests stay red after the reruns still end the run failed?

## Links

### In

- answers [plan-ci-flake-reruns](../decision/plan-ci-flake-reruns.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
