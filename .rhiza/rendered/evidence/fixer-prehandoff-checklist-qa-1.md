# fixer-prehandoff-checklist-qa-1

evidence · `01M788Q87H1GW2M37RADEYPZ4N` · actor [claude-code-22403536](../actor/claude-code-22403536.md)

## result

PASS

## statement

QA at 95fee8aca with a binary built from the head: `rhiza guide review` prints the Fixer section's checklist (reproduce the findings, `rhiza ship generate`, `rhiza render`, a Checks: block with exit codes); builtin review.toml's fixer and builtin ship-change.toml's ci-fix, gate-fix, gate-fix-merge, conflict-fix and conflict-fix-landed prompts all name the Checks: block. Caveat: in this repo `rhiza workflow describe ship-change` resolves to the repo layer .rhiza/workflows/ship-change.toml, which shadows the builtin and lacks the new text, so Rhiza's own ci-fix/gate-fix/conflict-fix fixers do not get it (filed as a discovered-from draft change set). go test -trimpath ./internal/guide ./internal/workflow pass.

## Links

### Out

- verifies [fixer-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
