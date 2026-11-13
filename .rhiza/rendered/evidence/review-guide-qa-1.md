# review-guide-qa-1

evidence · `01M79R58HKRJF0M91XT6A8X92E` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): ran rhiza-qa guide ls (exit 0): lists review and the nine role topics; rhiza-qa guide review prints the index naming all nine with skill and rhiza guide routes, scope routed to check-compliance; rhiza-qa guide scope exits 4 cleanly; guide install writes all ten review skills.

## Links

### Out

- verifies [review-guide](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
