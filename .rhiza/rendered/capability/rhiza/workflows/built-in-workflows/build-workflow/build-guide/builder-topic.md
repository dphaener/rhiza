# builder-topic

capability · `01M7974FWT13ZEJX4QJZQKJBFB` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a builder test that finds the six section headings in order, the shared done, test-scope, unstable-CI, command-budget, outside-work and never-list blocks expanded, the read-patterns-first and question-before-work rules, the acceptance commands, rhiza ship generate, rhiza render and the Checks: block above the one final trailer paragraph (the builder-prehandoff-checklist and guides-one-final-paragraph points, moved from build.md's tests), the actor-record git add line, and the RHIZA_ATTEMPT > 1 rule naming both /rhiza-build-retry and `rhiza guide build-retry`; a negative case asserts the builder text names no Go, make or repo-specific command.

## statement

`rhiza guide builder` gives a leaf's builder afk's builder prompt in Rhiza's words on six sections in order (Task, Expected outcome, Read first, Must do, Must not do, Context): the leaf's acceptance is the contract; read the slice, the leaf and the existing code and tests that follow the pattern before writing; raise a question the owner must answer before starting work, not midway; run the pre-handoff checklist and report a Checks: block; end with one commit carrying the trailers and the staged actor record, then session done as its own command; no stubs, TODOs, placeholders or new dependencies; and on attempt > 1 (RHIZA_ATTEMPT > 1) load build-retry first.

## Links

### Out

- depends-on [guide-shared-blocks](guide-shared-blocks.md)
- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [build-role-skills-installed](build-role-skills-installed.md)
- verifies [builder-topic-verified](../../../../../../evidence/builder-topic-verified.md)
- verifies [builder-topic-qa-1](../../../../../../evidence/builder-topic-qa-1.md)
- verifies [builder-topic-qa-2](../../../../../../evidence/builder-topic-qa-2.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
