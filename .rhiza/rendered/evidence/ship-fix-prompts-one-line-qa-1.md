# ship-fix-prompts-one-line-qa-1

evidence · `01M79R59B2M99WJ2954FZK9EKJ` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): ran rhiza-qa workflow describe ship-change and docs: every ci-fix, gate-fix, gate-fix-merge, conflict-fix, conflict-fix-landed and model-conflict-fix step is a one-line prompt naming its topic, docs.toml's matching the builtin; rhiza-qa workflow check ship-change and docs exit 0.

## Links

### Out

- verifies [ship-fix-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/ship-fix-prompts-one-line.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
