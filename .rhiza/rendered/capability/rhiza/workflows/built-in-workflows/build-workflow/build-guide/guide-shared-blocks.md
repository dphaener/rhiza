# guide-shared-blocks

capability · `01M7974FTA4EJ8NB0P869FGVT0` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with tests that (1) a topic naming a shared block prints its full text in place of the marker line, with no marker left in any listed text; (2) the expansion of a text naming an unknown block fails with an error naming the block (the load panics on it); (3) every listed text, expanded, stays within its cap; (4) the never-list block forbids sleeping or polling, stubs, TODOs and placeholders, new dependencies, rewriting history and chaining session done; and the budget, scoped-tests, unstable-CI, outside-work and done phrases build_test.go and budget_test.go check today are found in the shared files, each in exactly one file.

## statement

The guide keeps each shared block (session done rule, test scope, unstable CI, command budget, outside work, never-list) once as a markdown file under internal/guide/text/shared/, and a topic line `<!-- shared: <block> -->` is replaced by that block when the guide loads, so `rhiza guide <topic>`, the installed skills and the TopicCap all see the expanded text.

## Links

### Out

- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)
- derived-from [build-role-topics-interview](../../../../../../source/build-role-topics-interview.md)

### In

- depends-on [builder-topic](builder-topic.md)
- depends-on [verifier-topic](verifier-topic.md)
- depends-on [build-retry-topic](build-retry-topic.md)
- verifies [guide-shared-blocks-verified](../../../../../../evidence/guide-shared-blocks-verified.md)
- verifies [guide-shared-blocks-qa-1](../../../../../../evidence/guide-shared-blocks-qa-1.md)
- verifies [guide-shared-blocks-qa-2](../../../../../../evidence/guide-shared-blocks-qa-2.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
