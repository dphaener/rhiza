# check-quality-topic-verified

evidence · `01M79MY7S9315AC7K8NJZCNX5J` · actor [claude-code-d3b5ebd7](../actor/claude-code-d3b5ebd7.md)

## result

PASS

## statement

At ad9a381e0 re-ran the Checks: block: go test -trimpath ./internal/guide/... (ok, includes TestCheckQualityText and TestCheckQualityNamesNoRepoCommand), go test -trimpath ./internal/api -run Guide (ok), TestScripts/guide$ and TestScripts/guide-install$ (ok), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked change). Read the diff: internal/guide/text/check-quality.md has exactly the six sections in order; the built topic expands session-done, test-scope, command-budget, never-list (inside Must not do) and severity with no leftover markers; Task carries the isolation rule (knows nothing of how it was built, never rhiza run status / run show of the build run, its findings, run notes or transcripts, judges only slice contract, merge-base diff and Evidence); Read first discovers test commands from the contributor guide; Must do has the direct/integration/downstream impact radius, never silently fall back to the full suite, and the seven checklist items graded by severity. The test asserts all of these plus a negative case banning Go, make and repo-specific commands. Only the guide entry, list-order and ls-count tests, the txtar and the actor record changed; no stubs or TODOs.

## Links

### Out

- verifies [check-quality-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-quality-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
