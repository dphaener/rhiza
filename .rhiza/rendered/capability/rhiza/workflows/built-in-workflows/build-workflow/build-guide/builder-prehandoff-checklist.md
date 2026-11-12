# builder-prehandoff-checklist

capability · `01M76XY43WSHZDCT6WSCK7ATYJ` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide -run Build asserts build.md's builder section names the acceptance commands, rhiza ship generate, rhiza render and the Checks: block, and that the Checks: block sits above the one final trailer paragraph (guides-one-final-paragraph still holds); the build.toml builder prompt test in go test -trimpath ./internal/workflow asserts the prompt points at the checklist.

## statement

The build guide tells a builder, before its one leaf commit, to run the leaf's acceptance commands, then rhiza ship generate and rhiza render, fold every file they change into the commit, and list each command it ran with its exit code in a Checks: block of the commit body, above the trailer paragraph.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [builder-prehandoff-checklist-verified](../../../../../../evidence/builder-prehandoff-checklist-verified.md)
- verifies [builder-prehandoff-checklist-qa-1](../../../../../../evidence/builder-prehandoff-checklist-qa-1.md)
- verifies [builder-prehandoff-checklist-verified-2](../../../../../../evidence/builder-prehandoff-checklist-verified-2.md)
- verifies [pr-502-pre-handoff-checklist](../../../../../../evidence/pr-502-pre-handoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
