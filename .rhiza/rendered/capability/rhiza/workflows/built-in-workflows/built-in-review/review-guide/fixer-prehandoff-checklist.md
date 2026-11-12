# fixer-prehandoff-checklist

capability · `01M76XY4512EWBGYNTAAH4TQ4Y` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide -run Review asserts review.md's fixer section names rhiza ship generate, rhiza render and the Checks: block; go test -trimpath ./internal/workflow asserts the review.toml fixer prompt and ship-change.toml's ci-fix, gate-fix and conflict-fix prompts each name the Checks: block.

## statement

The review guide's fixer section and the ci-fix, gate-fix and conflict-fix prompts of the builtin ship-change workflow tell a fixer to run the commands that reproduce its findings, then rhiza ship generate and rhiza render, before its one commit, and to list each command with its exit code in a Checks: block of the commit body.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [fixer-prehandoff-checklist-verified](../../../../../../evidence/fixer-prehandoff-checklist-verified.md)
- verifies [fixer-prehandoff-checklist-qa-1](../../../../../../evidence/fixer-prehandoff-checklist-qa-1.md)
- verifies [fixer-prehandoff-checklist-verified-2](../../../../../../evidence/fixer-prehandoff-checklist-verified-2.md)
- verifies [pr-502-pre-handoff-checklist](../../../../../../evidence/pr-502-pre-handoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
