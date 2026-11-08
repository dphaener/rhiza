# ship-check-leaf-handle-renamed

decision · `01M6WK9BGGAN6KG625GFDBZJ7Y`

## reasoning

The leaf's statement says ship check 'never says a failure is unrelated to the plan' and its acceptance asserts no detail contains 'unrelated'; Decision plan-ci-flake-reruns chose a .rhiza-only fact over a cause claim because a check that reads .rhiza/ (Rhiza's own self-check) can fail on a plan's content, and Decision plan-ci-flake-source-summary-matches already corrected the Source to that design. Only the handle still names the rejected design, so a builder reading it, or plan-run-reruns-ci-flakes's dependency by handle, is pointed the wrong way. rhiza apply's set op takes a new handle on the same node, so the rename keeps the ID and every link; the dependency follows the node. The critic's suggested name describes exactly what the leaf adds: the rerun count and the .rhiza-only note. Changing the statement to match the handle instead would reintroduce the possibly false claim.

## statement

Rename leaf 01M6WJ6ZEEPYS69MVFGK79HHKN from ship-check-names-unrelated-failure to ship-check-failure-names-reruns-and-rhiza-only-diff (a set with a new handle on the same node, so its ID, fields, decided-by link to plan-ci-flake-reruns and plan-run-reruns-ci-flakes's depends-on link are kept). The leaf's statement and acceptance stay unchanged; any text that names the old handle is updated to the new one.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6whtv-r3-f2](../question/review-01m6whtv-r3-f2.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
