# verifier-topic

capability · `01M7974G189EFBV82T8PZ626EZ` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a verifier test that finds the six section headings in order, the shared done, test-scope, command-budget and never-list blocks expanded, the detached checkout, Never edit code, fakes only, the re-run of the Checks: block and acceptance commands with a missing block or a different exit code a finding naming the command (the verifier-reruns-reported-checks points, moved from build.md's tests), the scope and no-stub points, the Evidence apply example and the rhiza run findings heredoc, and the rule that each finding names a file, command or acceptance point; a negative case asserts it carries no builder commit or trailer instructions.

## statement

`rhiza guide verifier` gives a leaf's verifier afk's verifier prompt in Rhiza's words on the same six sections: in a detached checkout at the leaf's commit it re-runs every command of the Checks: block and every acceptance command, reads the diff and confirms each acceptance point is met and tested, nothing outside the leaf changed and no stubs or TODOs were left; it never edits code and uses fakes only; it ends with Evidence written as itself or with findings, each naming the file, command or acceptance point and what to change so a fixer can act without re-diagnosing, then session done as its own command.

## Links

### Out

- depends-on [guide-shared-blocks](guide-shared-blocks.md)
- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [build-role-skills-installed](build-role-skills-installed.md)
- verifies [verifier-topic-verified](../../../../../../evidence/verifier-topic-verified.md)
- verifies [verifier-topic-verified-2](../../../../../../evidence/verifier-topic-verified-2.md)
- verifies [verifier-topic-qa-2](../../../../../../evidence/verifier-topic-qa-2.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
