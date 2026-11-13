# review-01m79a41352xbz4fc2fjt5903j-r3-f4

question · `01M79BMA69T1ENTYGM5AMW29J0`

## statement

Plan review round 3 left finding f4 (contradiction on unstable-ci-guidance): Says the rhiza-review guide topic carries the unstable-CI reproduce-locally rule (TestReviewText finds it in review.md); the plan moves it to the shared unstable-CI block in fixer, ci-fix and gate-fix and makes review an index, without updating this node. The critic confirmed it: unstable-ci-guidance's acceptance finds the unstable-CI rule in the rhiza-review text via TestReviewText (review_test.go 'unstable-ci'); the plan moves it to the fixer/ci-fix/gate-fix topics' shared block and leaves this node naming the review topic. Major stands. How should it be settled?

## Links

### Out

- blocks [unstable-ci-guidance](../capability/rhiza/integrity/self-check/unstable-ci-guidance.md)

### In

- answers [unstable-ci-in-role-topics](../decision/unstable-ci-in-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
