# review-prompts-one-line-qa-1

evidence · `01M79R5930YR5WG5KAC1NAFG1D` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): ran rhiza-qa workflow describe review and docs-review: compliance, quality, qa and model checks only (no scope), each a one-line prompt naming its topic; critic and fixer one-line; docs-review fixer starts with the builtin line and keeps the humanize/doc-prose-humanized/doc-examples-plain-commands text; rhiza-qa workflow check review and docs-review exit 0.

## Links

### Out

- verifies [review-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-prompts-one-line.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
