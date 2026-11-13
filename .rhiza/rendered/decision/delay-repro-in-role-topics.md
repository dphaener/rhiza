# delay-repro-in-role-topics

decision · `01M79C6P7RASBNPZC64DXHM4VH`

## reasoning

Read the question, delay-repro-guidance (acceptance pins the rule to the review topic beside the unstable-CI rule via TestReviewText 'delay-repro'), review.md today (the delay rule is the line right after the unstable-CI rule), Decision review-role-topics-shape and leaf review-guide (review becomes an index; role topics carry rules through shared blocks), and the sibling rulings on this review round, above all unstable-ci-in-role-topics, which moved the unstable-CI rule into the shared unstable-CI block expanded into fixer, ci-fix and gate-fix. The node is stale, not the plan. Putting the delay rule in that same block keeps it beside the unstable-CI rule as its acceptance already requires, needs no new block, and makes the three leaves whose acceptances already require the unstable-CI block carry it without new leaf work. I weighed also expanding it into check-quality and check-qa: checks run scoped tests, but a check that meets a flake records it as a finding and does not chase it; reproducing a race is fixer work, and the checks' host-safety is already covered by the shared command-budget block (no sleep or poll, about four minutes per command, background once). conflict-fix resolves merges, not reds, so it is left out as in the unstable-CI ruling. No CLI command, flag, --json shape, rhiza.* setting or .rhiza/ format changes, so no breaking note.

## statement

delay-repro-guidance is updated in this change set to follow the split, not by keeping the rule in the review index: AGENTS.md and the rhiza-build topic keep the delay-injection rule unchanged (TestClaudeMDScopedTests and TestBuildText as today), and on the review side the rule (inject a delay at the suspected race; never busy loops, CPU stress or parallel stress runs on the shared host) joins the shared unstable-CI block, beside the unstable-CI reproduce-locally rule, so it is expanded into the fixer, ci-fix and gate-fix topics; the node's statement names builders, verifiers and those three fixer roles instead of rhiza-review's checks and fixers, and its acceptance finds the rule in the expanded fixer, ci-fix and gate-fix texts through the guide package tests instead of TestReviewText on review.md.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f8](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f8.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
