# command-budget-in-role-topics

decision · `01M79BSSTB8CZXWB8WVSJVBB5Q`

## reasoning

Read the question, command-budget-guide (acceptance pins the rule to build.md and review.md; budget_test.go loops over build and review), Decision review-role-topics-shape (review becomes a short index; the role topics reuse the shared blocks, command-budget among them), leaf review-role-skills-installed (review.md's old phrases checked by review_test.go and budget_test.go must be found in the role topics) and the role topic leaves, whose acceptances already require the shared command-budget block expanded. The node is stale, not the plan: keeping the rule in the index would duplicate it and contradict the index shape the arbiter already chose. This matches the ruling on the sibling finding for guides-scoped-tests (scoped-tests-in-role-topics), so the two shared-block rules are settled the same way. No CLI, flag, --json, setting or .rhiza/ format changes, so no breaking note.

## statement

command-budget-guide is updated in this change set to follow the split, not by keeping the budget rule in the review index: the build topic keeps the rule for builders and verifiers (build.md unchanged), and on the review side the rule lives in the shared command-budget block expanded into every review role topic (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix), not in review.md. Its statement names the build topic and the review role topics, and its acceptance becomes: go test -trimpath ./internal/guide/... asserts build.md and each expanded review role topic hold the budget rule (no sleep or poll, about four minutes per command, background once and wait once) and that the rule names no Go, make or Rhiza-repo command; budget_test.go loops over build plus the nine role topics instead of build and review.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f3](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f3.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
