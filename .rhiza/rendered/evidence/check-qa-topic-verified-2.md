# check-qa-topic-verified-2

evidence · `01M79S5ZF9R7SF55EQW6NRTX3C` · actor [claude-code-1a97b1e8](../actor/claude-code-1a97b1e8.md)

## result

PASS

## statement

Attempt 2 at f62ce2fb9: re-ran go test -trimpath -count=1 ./internal/guide/... (exit 0; TestCheckQAText now also asserts the shared test-scope block, TestCheckQANamesNoRepoCommand passes), built rhiza to $TMPDIR and ran `guide check-qa | grep -q "only the tests covering"` (exit 0), `ship generate` (exit 0, 0 files changed), `render` (exit 0, no tracked changes) and make self-check (exit 0). Diff touches only check-qa.md, check_qa_test.go, rendered actors/index and the builder's actor record. The six sections stay in order, isolation, clean-state, scratch path and never-install (still under Must do), happy/failure, bounded exploration, broken-only and the <leaf>-qa-<cycle> verifies Evidence example are all still present and tested. No stubs or TODOs.

## Links

### Out

- verifies [check-qa-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-qa-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
