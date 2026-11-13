# verifier-reruns-reported-checks

capability · `01M76XY462ZXBJ1BQB0RJ5FA66` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide -run TestVerifierText asserts verifier.md says to re-run every command in the Checks: block and every acceptance command, and that a missing Checks: block or a different exit code is a finding naming the command.

## statement

The verifier guide topic (rhiza guide verifier, so the rhiza-verifier skill) tells a verifier to re-run every command in the leaf commit's Checks: block and every acceptance command, and to record a finding naming the command when one exits differently from what was reported or the commit has no Checks: block.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [verifier-reruns-reported-checks-verified](../../../../../../evidence/verifier-reruns-reported-checks-verified.md)
- verifies [verifier-reruns-reported-checks-qa-1](../../../../../../evidence/verifier-reruns-reported-checks-qa-1.md)
- verifies [pr-502-pre-handoff-checklist](../../../../../../evidence/pr-502-pre-handoff-checklist.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
