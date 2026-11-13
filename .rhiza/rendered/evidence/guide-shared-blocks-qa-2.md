# guide-shared-blocks-qa-2

evidence · `01M79C8NCXMSE2V2CS318Q1T56` · actor [claude-code-9b656f6b](../actor/claude-code-9b656f6b.md)

## result

PASS

## statement

`rhiza guide builder` and `rhiza guide verifier` (built at HEAD) print the shared test-scope, unstable-ci, command-budget, outside-work, never-list and session-done text expanded, with zero `<!--` markers left in output or in the installed SKILL.md files. go test -trimpath ./internal/guide/... ok

## Links

### Out

- verifies [guide-shared-blocks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/guide-shared-blocks.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
