# conflict-fix-topic-verified

evidence · `01M79P85R85XKJHZWRFJJBXBKM` · actor [claude-code-407a92e0](../actor/claude-code-407a92e0.md)

## result

PASS

## statement

At 91c8fa06d re-ran the Checks: block: go test -trimpath ./internal/guide/... (incl. TestConflictFixText and the negative TestConflictFixNamesNoRepoCommand), go test -trimpath ./internal/api -run Guide, go test -trimpath ./cmd/rhiza -run TestScripts/guide, rhiza ship generate, rhiza render: all exit 0. A scratch build's `rhiza guide conflict-fix` prints the six sections in order with the session-done, test-scope, command-budget and never-list blocks expanded; the code part (rhiza ship check --json merge_conflict, ship generate, render, git add -A -- .rhiza/rendered .rhiza/log, a Checks: block, Change:/Actor:/Review: trailers), the model part (rhiza ship check --model-conflict --json, rhiza set <node> --rebase, Model-Conflict: trailer), the never-end-with-the-merge-in-progress rule, the one-blocking-question rule naming both change sets, and never push, rebase or force. The diff touches only the guide topic, its registration, list-count tests and the builder's actor record; no TODOs.

## Links

### Out

- verifies [conflict-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/conflict-fix-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
