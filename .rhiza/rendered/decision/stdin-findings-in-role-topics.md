# stdin-findings-in-role-topics

decision · `01M79C1ASS1P13DDNGR36CR9HS`

## reasoning

The finding and the critic agree the node is stale: it pins the heredocs to the review guide via review_test.go, while review-guide makes review an index and the role topics carry the work. The sibling rulings on the same round (scoped-tests-in-role-topics, command-budget-in-role-topics, done-own-command-in-role-topics, unstable-ci-in-role-topics) all settle this shape the same way, by updating the node to name the role topics rather than keeping text in the index, so this follows them. check-compliance-topic and critic-topic acceptances already name the findings and verdict heredocs, but check-quality, check-qa and check-model do not, though each check must record findings (a check with none fails); adding the heredoc there keeps the node's coverage for every check rather than silently dropping it. Fixers record no findings or verdicts, so they are left out; the critic's negative case (no findings example) stays consistent.

## statement

guides-stdin-findings is updated in this change set to follow the split, not by keeping the heredocs in the review index: the build and plan-review guides keep their heredocs unchanged (build_test.go and plan_test.go as today), and on the review side the `rhiza run findings --check` heredoc lives in each check topic (check-compliance, check-quality, check-qa, check-model) and the `rhiza run verdict --check` heredoc in the critic topic, each with its existing flags and JSON shape ({"findings":[]} when clean); the review index shows neither. Its statement becomes: the build guide, the plan-review guide and the four check topics and the critic topic show rhiza run findings and rhiza run verdict fed from stdin by a heredoc. Its acceptance becomes: tests in internal/guide assert build.md and plan-review.md as today, assert each check topic shows the `- <<'EOF'` form of rhiza run findings and the critic topic that of rhiza run verdict, and that no `<file>` form of either command remains in any of them; go test -trimpath ./internal/guide/... passes. check-quality, check-qa and check-model, whose acceptances do not name the findings heredoc today, gain it so every check can record its findings.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f6](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f6.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
