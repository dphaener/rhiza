# review-01m79a41352xbz4fc2fjt5903j-r3-f8

question · `01M79BMAFG4MXMP85H1AR405ZF`

## statement

Plan review round 3 left finding f8 (contradiction on delay-repro-guidance): Says the rhiza-review guide topic tells checks and fixers to reproduce a timing flake by injecting a delay; the plan empties review.md into an index and does not say which role topic (fixer, ci-fix, gate-fix, check-quality) carries this rule, nor update this node. The critic confirmed it: delay-repro-guidance's acceptance finds the delay-injection rule in the review topic via TestReviewText ('delay-repro'); review.md becomes an index, and no plan leaf's acceptance (fixer, ci-fix, gate-fix, check-quality) names the delay rule, so the rule is either lost or the landed acceptance fails. Major stands. How should it be settled?

## Links

### Out

- blocks [delay-repro-guidance](../capability/rhiza/integrity/self-check/delay-repro-guidance.md)

### In

- answers [delay-repro-in-role-topics](../decision/delay-repro-in-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
