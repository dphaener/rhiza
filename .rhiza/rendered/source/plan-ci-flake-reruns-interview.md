# plan-ci-flake-reruns-interview

source · `01M6WJ6B8CAWDTWP2YT6NMESGW`

## quotes

```json
[]
```

## raw_sha256

dc48b438eb47615f227d087a39da0a5fe8e8ad9e0546bc6029ec4e0456c33f73

## session

5d744c51-b515-4500-b578-01a1b91df18d

## summary

Autonomous plan run for change set 01M6WHT. The planner found that Nebari plan run 01M6VY5APJP9 (PR #78) failed at rebased-green twice, on two different rebased heads, each on a different flaky RSpec system spec, after ship check's existing once-per-head rerun of the failed Actions jobs had already fired and failed again (the ship record shows rerun_head and reruns for e893ad0); a manual third attempt passed. The planner asked four questions (which intent the change set serves, how the plan workflow gets an extra rerun, how the failure is named as unrelated to the plan, and whether a .rhiza-only PR that stays red still fails). With no owner, the arbiter answered all four with Decision plan-ci-flake-reruns: serve m27-nebari-on-rhiza; rhiza ship check gains --reruns <n> (default 1) with a per-head count in the ship record, and plan.toml passes --reruns 2 to its four ship check steps; once the reruns are spent, the error names the failed checks, adds ' after <n> reruns' when reruns were spent, and the parsed failing tests, and, when the PR's diff from its merge base touches only .rhiza/, appends that the PR changes only .rhiza/ (a fact, not a claim about the failure's cause); the run still fails on a persistent red. No owner words were given.

## title

Plan workflow CI flake reruns, autonomous plan interview

## Links

### In

- derived-from [plan-ci-flake-reruns](../decision/plan-ci-flake-reruns.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
