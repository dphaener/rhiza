# builder-prehandoff-checklist

capability · `01M76XY43WSHZDCT6WSCK7ATYJ` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide -run TestBuilderText asserts internal/guide/text/builder.md's Must do section names the acceptance commands, rhiza ship generate, rhiza render and the Checks: block with exit codes, and that the Checks: block sits above the one final trailer paragraph (guides-one-final-paragraph still holds).

## statement

The builder guide topic (rhiza guide builder) tells a builder, before its one leaf commit, to run the leaf's acceptance commands, then rhiza ship generate and rhiza render, fold every file they change into the commit, and list each command it ran with its exit code in a Checks: block of the commit body, above the trailer paragraph. The one-line builder prompt only points at the topic.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [builder-prehandoff-checklist-verified](../../../../../../evidence/builder-prehandoff-checklist-verified.md)
- verifies [builder-prehandoff-checklist-qa-1](../../../../../../evidence/builder-prehandoff-checklist-qa-1.md)
- verifies [builder-prehandoff-checklist-verified-2](../../../../../../evidence/builder-prehandoff-checklist-verified-2.md)
- verifies [pr-502-pre-handoff-checklist](../../../../../../evidence/pr-502-pre-handoff-checklist.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
