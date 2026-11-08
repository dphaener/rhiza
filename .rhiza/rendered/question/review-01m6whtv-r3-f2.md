# review-01m6whtv-r3-f2

question · `01M6WK5KK0ZXXJN489FHC7KJ27`

## statement

Plan review round 3 left finding f2 (contradiction on ship-check-names-unrelated-failure): The leaf's handle says ship check names an unrelated failure, but its statement and acceptance forbid that: ship check 'never says a failure is unrelated to the plan', and no detail may contain 'unrelated'. A builder reading the handle, and plan-run-reruns-ci-flakes's depends-on that names it, gets the rejected design. Rename it to match the statement, for example ship-check-failure-names-reruns-and-rhiza-only-diff. The critic confirmed it: Checked: the leaf's handle is ship-check-names-unrelated-failure while its statement says ship check 'never says a failure is unrelated to the plan' and its acceptance asserts no detail contains 'unrelated'. plan-run-reruns-ci-flakes depends on it by that handle. The handle names the design rounds 1-2 rejected and contradicts the leaf's own statement; renaming it to describe the rerun count and .rhiza-only note fixes it. How should it be settled?

## Links

### Out

- blocks [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)

### In

- answers [ship-check-leaf-handle-renamed](../decision/ship-check-leaf-handle-renamed.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
