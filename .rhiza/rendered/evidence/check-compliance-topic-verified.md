# check-compliance-topic-verified

evidence · `01M79MQ35V5TT48SHGQSZJRK97` · actor [claude-code-b1e9e133](../actor/claude-code-b1e9e133.md)

## result

PASS

## statement

At cc7f78ac9 re-ran the Checks: block: go test -trimpath ./internal/guide/... (exit 0), go test -trimpath ./internal/api -run Guide (exit 0), go test -trimpath ./cmd/rhiza -run TestScripts/guide$ and TestScripts/guide-install$ (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0). internal/guide/check_compliance_test.go asserts the six sections in order, the session-done, test-scope, command-budget, never-list and severity blocks expanded (never-list inside Must not do), the isolation rule in Task, the under-built/scope-creep/unaccounted-files criteria with severities, the one-line-with-evidence rule and the rhiza run findings --check heredoc with {"findings":[]}; a negative test bans Go, make, repo paths, Checks:, trailer and commit text. rhiza guide check-compliance output read and matches. Diff touches only the new topic, its registry entry, guide list tests and the builder's actor record; no TODOs.

## Links

### Out

- verifies [check-compliance-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-compliance-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
