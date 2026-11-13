# check-model-topic-verified

evidence · `01M79NBAZ1GV826FDD4AWCMA8H` · actor [claude-code-64a16129](../actor/claude-code-64a16129.md)

## result

PASS

## statement

At fdbc00896 re-ran the Checks block: go test -trimpath ./internal/guide/... (TestCheckModelText and TestCheckModelNamesNoRepoCommand PASS with -count=1 -v), go test -trimpath ./internal/api -run Guide, go test -trimpath ./cmd/rhiza -run TestScripts/guide$ and TestScripts/guide-install$, rhiza ship generate, rhiza render: all exit 0 with no tracked file changed. The diff adds internal/guide/text/check-model.md, its topic entry, the list-order and guide.txtar counts, and the test. The test asserts the six sections in order, the session-done, command-budget, outside-work, never-list (inside Must not do) and severity blocks expanded, the isolation rule in Task (no run status/run show of the build run, its findings, run notes or transcripts; judges the slice contract, merge-base diff and model Evidence), code-contradicts-node major and missing-constraint minor, the discovered-from change new/link commands, and a negative case banning Go, make and repo-specific commands. Rendered rhiza guide check-model output reads correctly. Nothing outside the leaf changed; no stubs or TODOs.

## Links

### Out

- verifies [check-model-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-model-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
