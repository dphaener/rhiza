# review-role-skills-installed-qa-1

evidence · `01M79R58VC6KHDAA8K92XHP0DT` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): ran rhiza-qa guide install --agent claude into an empty scratch HOME: 19 written incl. the nine role skills with their topic description in frontmatter; rerun 0 written, 19 unchanged; over an install by the old rhiza: 11 written, 8 unchanged, no --force needed; --agent bogus exits 2 with a next: hint.

## Links

### Out

- verifies [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
