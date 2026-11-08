# plan-ci-flake-source-summary-matches

decision · `01M6WK7C2W9WTGM1MWZE4PEG5G`

## reasoning

The Source's summary is the planner's own paraphrase of the arbiter's answer, not owner words (quotes is empty; 'No owner words were given'), so correcting it invents nothing. The Decision it backs explicitly rejects claiming the cause: a check that reads .rhiza/, such as Rhiza's own self-check, can fail on a plan's content, so 'unrelated' could be false; the leaf forbids 'unrelated' in any detail. The Decision and leaf are the reviewed, authoritative design, and the change set's intent ('naming the test as unrelated') is served by naming the failing tests and stating the .rhiza-only fact. Changing the Decision to match the Source would reintroduce the rejected, possibly false claim, so the Source must follow the Decision so provenance readers see the chosen design.

## statement

Correct Source plan-ci-flake-reruns-interview's summary to match Decision plan-ci-flake-reruns: replace 'says the failure is unrelated to the plan' with wording that, when the PR's diff from its merge base touches only .rhiza/, the error appends that the PR changes only .rhiza/ (a fact, not a claim about the cause), and adds ' after <n> reruns' when reruns were spent. Decision plan-ci-flake-reruns and leaf ship-check-names-unrelated-failure stay as they are; the Source summary is the only thing that changes, and it adds no owner words since none were given.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6whtv-r3-f1](../question/review-01m6whtv-r3-f1.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
