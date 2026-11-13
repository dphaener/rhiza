# done-own-command-in-role-topics

decision · `01M79BYWK41YAXH8WS8EWSATG5`

## reasoning

Read the question, done-own-command-guide (acceptance pins the rule to build.md and review.md; review_test.go 'done' reads review.md), Decision review-role-topics-shape and leaf review-role-skills-installed (review.md becomes a short index; its old phrases checked by review_test.go must be found in the role topics), and the role topic leaves, whose acceptances already require the shared done block expanded. The node is stale, not the plan: keeping the rule in the index would duplicate it and contradict the index shape already chosen. This matches the rulings on the sibling findings (scoped-tests-in-role-topics, command-budget-in-role-topics, unstable-ci-in-role-topics), so every shared-block rule is settled the same way. All nine topics end with rhiza session done, so all nine carry the block, and no role the node covered loses the rule. The workflow half stays: the one-line prompts drop session-done text, which the existing prompt test already allows. No CLI command, flag, --json shape, rhiza.* setting or .rhiza/ format changes, so no breaking note.

## statement

done-own-command-guide is updated in this change set to follow the split, not by keeping the rule in the review index: build.md keeps the own-command rule for builders and verifiers unchanged, and on the review side the rule lives in the shared done block expanded into every review role topic (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix), not in review.md. Its statement names the build topic, the review role topics and the builtin prompts, and its acceptance becomes: go test -trimpath ./internal/guide/... asserts build.md and each of the nine expanded review role topics hold the own-command rule (run rhiza session done as its own command, only after the recording command succeeded, never chained with ; or &&); go test -trimpath ./internal/workflow asserts no prompt in internal/workflow/builtin/*.toml contains '; rhiza session done' or '&& rhiza session done' and every prompt that still names rhiza session done says it is the very last action (one-line /rhiza-<role> prompts that no longer name it pass).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f5](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f5.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
