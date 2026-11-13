# scoped-tests-in-role-topics

decision · `01M79BR06PH8ZCHF6PKW34TC1H`

## reasoning

The finding and the critic agree the node is stale: it pins the rule to review.md via TestReviewText, while review-guide turns review.md into an index and review-role-skills-installed requires review.md's old phrases checked by review_test.go to be found in the role topics instead. So the node must name the role topics. I weighed narrowing the node to only the topics whose plan acceptances already list the test-scope block (check-compliance, check-quality, fixer, ci-fix, gate-fix, conflict-fix), but that would silently drop the rule for check-qa, check-model and critics, which guides-scoped-tests (decided by ci-runs-full-suite and scoped-tests-serve-m25-generic-rule) covers today; this change set restructures the review guide and is not meant to weaken that rule. Including the shared block in those three topics is a small, consistent addition (the block is shared and command-free) and keeps the node's coverage while moving its home from review.md to the role topics.

## statement

guides-scoped-tests is updated to follow the split: the build topic keeps the rule for builders and verifiers (TestBuildText on build.md unchanged), and on the review side the rule lives in the shared test-scope block expanded into every review role topic (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix), not in the review index. Its acceptance becomes: the build test asserts build.md states the rule, and the internal/guide tests assert each of the nine expanded role topics contains the test-scope block and none of build.md or the role topics contains go test, make check, make test or testscript; go test -trimpath ./internal/guide/... passes. The check-qa, check-model and critic topics, whose acceptances do not list test-scope today, gain the block so no role the node covered loses the rule.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f2](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f2.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
