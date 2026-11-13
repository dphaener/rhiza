# gate-fix-topic-qa-1

evidence · `01M79R583THJ404W6E068CHRME` · actor [claude-code-0d4e7004](../actor/claude-code-0d4e7004.md)

## result

PASS

## statement

QA at HEAD 1eaad5503 with a scratch build (go build -trimpath -o <scratchpad>/rhiza-qa ./cmd/rhiza, not installed): ran rhiza-qa guide gate-fix: exit 0, renders Task, Expected outcome, Read first, Must do, Must not do, Context in order with shared blocks expanded (no unexpanded placeholders); it is listed by rhiza-qa guide ls and installed as skills/rhiza-gate-fix/SKILL.md by rhiza-qa guide install --agent claude into a scratch HOME (rerun: 0 written, 19 unchanged). Failure: rhiza-qa guide nosuch exits 4 with a next: hint. Text carries the Full log rule and never rerunning rhiza ship gate.

## Links

### Out

- verifies [gate-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/gate-fix-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
