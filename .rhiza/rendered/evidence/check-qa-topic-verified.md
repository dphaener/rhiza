# check-qa-topic-verified

evidence · `01M79N5BXPXJHSVVTVAB4VNE7W` · actor [claude-code-ddad1012](../actor/claude-code-ddad1012.md)

## result

PASS

## statement

At af9fc072e re-ran the Checks block: go test -trimpath ./internal/guide/... (also -count=1: TestCheckQAText, TestCheckQANamesNoRepoCommand, TestListOrder, shared tests PASS), go test -trimpath ./internal/api -run Guide, TestScripts/guide$ and TestScripts/guide-install$, rhiza ship generate and rhiza render (scratch build of this commit): all exit 0, no tracked changes. Diff read: internal/guide/text/check-qa.md has the six sections in order, expands session-done, command-budget, never-list (inside Must not do) and severity, carries the isolation rule in Task, clean state, scratch path outside $PATH and never-install, happy and failure scenarios, bounded pairs then edge cases with no new test files, broken-behaviour-only, and the Evidence apply example with <leaf>-qa-<cycle> and verifies; check_qa_test.go asserts each, and a negative test asserts no Go, make or repo command. Only the guide registry, list tests, guide.txtar count and the leaf's own log/render files changed; no stubs or TODOs.

## Links

### Out

- verifies [check-qa-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-qa-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
