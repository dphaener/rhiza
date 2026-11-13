# one-final-paragraph-in-fix-topics

decision · `01M79C47H71KK6S71389S322M1`

## reasoning

Read the question, guides-one-final-paragraph (its test TestReviewFixerTrailersInOneFinalParagraph reads review.md, where the rule sits in the Fixer section today), review-guide and Decision review-role-topics-shape (review becomes an index; role topics carry the work), and the fixer-topic and conflict-fix-topic acceptances, which already name the one final trailer paragraph and the Review:/Model-Conflict: trailers. The node is stale, not the plan: keeping the rule in the index would contradict the chosen index shape. The sibling rulings on the same round (scoped-tests-in-role-topics, command-budget-in-role-topics, done-own-command-in-role-topics, unstable-ci-in-role-topics, stdin-findings-in-role-topics) all update the node to name the role topics, so this follows them. Only fixers commit on the review side, and ci-fix, gate-fix and conflict-fix commit too, so all four carry the rule and no committing role the node covered loses it; checks and critics write no commits and are left out. No CLI command, flag, --json shape, rhiza.* setting or .rhiza/ format changes, so no breaking note.

## statement

guides-one-final-paragraph is updated in this change set to follow the split, not by keeping a fixer instruction in the review index: build.md keeps the builder example unchanged, and on the review side the one-final-paragraph trailer rule (every trailer, Co-Authored-By: included, in one final paragraph, no blank line) lives in each topic whose role commits: fixer, ci-fix, gate-fix and conflict-fix, not in review.md. Its statement becomes: the build guide's builder example and the fixer, ci-fix, gate-fix and conflict-fix topics show every trailer, Co-Authored-By included, in one final paragraph. Its acceptance becomes: go test -trimpath ./internal/guide/... asserts build.md as today and each of the four fix topics holds Co-Authored-By:, one final paragraph and no blank line (TestReviewFixerTrailersInOneFinalParagraph moves from review.md to those topics), and fails when any of them puts Co-Authored-By in a paragraph apart from Change:, Actor:, Leaf:, Review: or Model-Conflict:.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f7](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f7.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
