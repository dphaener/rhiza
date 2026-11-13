# severity-shared-block-verified

evidence · `01M79J04N2GYTG2DWQJKEVHKVE` · actor [claude-code-6d8916eb](../actor/claude-code-6d8916eb.md)

## result

PASS

## statement

At 4af78f087, re-ran the commit's Checks: go test -trimpath ./internal/guide/... exit 0, go test -trimpath ./internal/api -run Guide exit 0, rhiza ship generate exit 0 (0 files changed), rhiza render exit 0 (tree clean after). Read the diff: severity.md holds the four severities and a mapping naming compliance (now with work outside the intent major), quality, qa and model, no scope check, plus plan review's kinds, and no critic outcomes or Go/make/Rhiza-repo commands. plan-review.md and review.md include it with <!-- shared: severity -->. A scratch build's `guide plan-review` shows the expanded block. TestSeverityBlockOnceInEachTopic checks that plan-review, review and every check-*, critic and fixer topic hold the block exactly once. TestSeverityBlockContent checks the severities, the mapping, the kinds, no scope, no critic outcomes and no banned commands, and that confirmed/downgraded/dropped sit in the critic sections. The check-*, critic and fixer topics come from dependent leaves, and the tests already cover them. Nothing outside the leaf changed besides the builder's actor record. No stubs.

## Links

### Out

- verifies [severity-shared-block](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-shared-block.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
