# severity-shared-block-qa-1

evidence · `01M79R59J6QE9F7H76ST82H5K6` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): read the Severities block in rhiza-qa guide output for check-compliance, check-quality, check-qa, check-model, critic, fixer and plan-review: the same block, four severities, mapping names compliance, quality, qa and model and no scope check.

## Links

### Out

- verifies [severity-shared-block](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-shared-block.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
