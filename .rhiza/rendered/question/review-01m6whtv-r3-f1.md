# review-01m6whtv-r3-f1

question · `01M6WK5KH8ZSWGNGMXXYGT0X0P`

## statement

Plan review round 3 left finding f1 (contradiction on plan-ci-flake-reruns-interview): The Source's summary says that once the reruns are spent the error 'says the failure is unrelated to the plan' when the diff touches only .rhiza/. Decision plan-ci-flake-reruns, which is derived from this Source, and leaf ship-check-names-unrelated-failure say the opposite: the note is only '; the PR changes only .rhiza/', it never claims the cause, and no detail may contain 'unrelated'. The summary has to match the Decision it backs. The critic confirmed it: Checked: Source plan-ci-flake-reruns-interview's summary says that with a .rhiza-only diff the error 'says the failure is unrelated to the plan', and attributes that to Decision plan-ci-flake-reruns. The Decision (derived-from this Source) now says the note is only that the PR changes only .rhiza/, 'a fact and not a claim about the cause', and leaf ship-check-names-unrelated-failure forbids 'unrelated' in any detail. The Source misreports the Decision it backs, so a reader following provenance gets the rejected design; the summary must match (no owner words are involved, so correcting it invents nothing). How should it be settled?

## Links

### In

- answers [plan-ci-flake-source-summary-matches](../decision/plan-ci-flake-source-summary-matches.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
