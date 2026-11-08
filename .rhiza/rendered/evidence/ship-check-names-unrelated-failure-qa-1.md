# ship-check-names-unrelated-failure-qa-1

evidence · `01M6XJKBP2DAATHFH107E4ERHR` · actor [claude-code-14b9e333](../actor/claude-code-14b9e333.md)

## result

PASS

## statement

QA cycle 1, by hand with the built binary and a fake gh: a head whose diff from origin/main touches only .rhiza/x/note, red after its two reruns, failed with 'PR #78 at 378f1f9: failed: test (4/5) after 2 reruns; the PR changes only .rhiza/'. A head that also touches app.go, red after one rerun with a failed-run log naming an RSpec test, failed with 'PR #78 at c06657e: failed: test (4/5) after 1 reruns; failing tests: rspec ./spec/system/garden_location_spec.rb:11' and no .rhiza/ note. A red check with no Actions run (no rerun spent) printed today's 'PR #78 at 0fb5524: failed: test (4/5)'. No detail said unrelated, flaky or on main.

## Links

### Out

- verifies [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
