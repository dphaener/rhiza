# review-01m79a41352xbz4fc2fjt5903j-r3-f2

question · `01M79BMA240JCRD8ME3SZVEKRQ`

## statement

Plan review round 3 left finding f2 (contradiction on guides-scoped-tests): Says the review guide topic (and the rhiza-review skill) states the scoped-tests rule for checks, critics and fixers, asserted by TestReviewText on review.md; the plan turns review into a short index and moves the rule into the check-*/critic/fixer topics via the shared test-scope block, but leaves this node unchanged. Update it to name the role topics. The critic confirmed it: guides-scoped-tests' acceptance has TestReviewText assert review.md states the rule for checks, critics and fixers (review_test.go 'scoped-tests' phrases); the plan cuts review.md to an index and leaves this node unchanged. Major stands. How should it be settled?

## Links

### Out

- blocks [guides-scoped-tests](../capability/rhiza/guide/guides-scoped-tests.md)

### In

- answers [scoped-tests-in-role-topics](../decision/scoped-tests-in-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
