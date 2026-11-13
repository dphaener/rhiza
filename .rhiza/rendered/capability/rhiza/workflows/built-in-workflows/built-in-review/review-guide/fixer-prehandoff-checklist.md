# fixer-prehandoff-checklist

capability · `01M76XY4512EWBGYNTAAH4TQ4Y` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... asserts the expanded fixer, ci-fix, gate-fix and conflict-fix topics each name rhiza ship generate, rhiza render and the Checks: block; go test -trimpath ./internal/workflow/... asserts the review.toml fixer prompt names the fixer topic and ship-change.toml's ci-fix, gate-fix and conflict-fix prompts name the ci-fix, gate-fix and conflict-fix topics.

## statement

The fixer, ci-fix, gate-fix and conflict-fix topics tell a fixer to run the commands that reproduce its findings, then rhiza ship generate and rhiza render, before its one commit, and to list each command with its exit code in a Checks: block of the commit body; the review workflow fixer prompt and the builtin ship-change workflow ci-fix, gate-fix and conflict-fix prompts reach them by naming their topic.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [fixer-prehandoff-checklist-verified](../../../../../../evidence/fixer-prehandoff-checklist-verified.md)
- verifies [fixer-prehandoff-checklist-qa-1](../../../../../../evidence/fixer-prehandoff-checklist-qa-1.md)
- verifies [fixer-prehandoff-checklist-verified-2](../../../../../../evidence/fixer-prehandoff-checklist-verified-2.md)
- verifies [pr-502-pre-handoff-checklist](../../../../../../evidence/pr-502-pre-handoff-checklist.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
