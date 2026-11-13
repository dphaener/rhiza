# guides-one-final-paragraph

capability · `01M62RE2VD8QDWJ0BQN4Y56F9K` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes: TestBuilderText asserts the builder example keeps every trailer in one final paragraph, TestReviewFixerTrailersInOneFinalParagraph asserts each of the fixer, ci-fix, gate-fix and conflict-fix topics holds Co-Authored-By:, one final paragraph and no blank line, and TestFixerText, TestCIFixText, TestGateFixText and TestConflictFixText fail when their example puts Co-Authored-By in a paragraph apart from Change:, Actor:, Leaf: or Review:.

## statement

The builder guide topic's example commit and the fixer, ci-fix, gate-fix and conflict-fix topics, the review roles that commit, show every trailer, Co-Authored-By included, in one final paragraph with no blank line; the review index holds no fixer instruction.

## Links

### In

- blocks [review-01m79a41352xbz4fc2fjt5903j-r3-f7](../../../question/review-01m79a41352xbz4fc2fjt5903j-r3-f7.md)
- verifies [guides-one-final-paragraph-verified](../../../evidence/guides-one-final-paragraph-verified.md)
- verifies [guides-one-final-paragraph-qa-1](../../../evidence/guides-one-final-paragraph-qa-1.md)
- verifies [pr-61-leaf-trailer-split-paragraphs](../../../evidence/pr-61-leaf-trailer-split-paragraphs.md)
- verifies [pr-537-review-role-topics](../../../evidence/pr-537-review-role-topics.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
