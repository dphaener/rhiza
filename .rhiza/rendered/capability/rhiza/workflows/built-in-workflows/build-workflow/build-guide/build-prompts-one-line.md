# build-prompts-one-line

capability · `01M7974GD79WS6HEB3RMEMQ7W9` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/workflow/... passes with the builtin builder and verifier prompt tests updated to the one-line prompts (no checklist, Checks:, session done or test-scope text left in them) and a test that review.toml's verifier prompt equals build.toml's; a test reads .rhiza/workflows/docs-build.toml and .rhiza/workflows/docs-review.toml and asserts each one's verifier prompt equals the builtin one and docs-build's builder prompt starts with the builtin line and keeps the humanize addendum; go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-evidence-commit|docs-commands-exist)' passes; the describe-golden script is regenerated only if its output changed; make self-check passes.

## statement

The builtin build workflow's builder prompt is the one line `/rhiza-builder Build leaf {leaf} of change set {change} ({title}), attempt {attempt}: follow rhiza guide builder.` and its verifier prompt, and review.toml's leaf verifier prompt, the one line `/rhiza-verifier Verify leaf {leaf} of change set {change} ({title}) at commit $RHIZA_COMMIT: follow rhiza guide verifier.`; Rhiza's .rhiza/workflows/docs-build.toml uses the builder line plus its unchanged doc-prose addendum and the verifier line byte for byte, and .rhiza/workflows/docs-review.toml's leaf verifier prompt is the verifier line byte for byte.

## Links

### Out

- depends-on [build-role-skills-installed](build-role-skills-installed.md)
- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- verifies [build-prompts-one-line-verified](../../../../../../evidence/build-prompts-one-line-verified.md)
- verifies [build-prompts-one-line-qa-1](../../../../../../evidence/build-prompts-one-line-qa-1.md)
- verifies [build-prompts-one-line-verified-3](../../../../../../evidence/build-prompts-one-line-verified-3.md)
- verifies [build-prompts-one-line-qa-2](../../../../../../evidence/build-prompts-one-line-qa-2.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
