# build-prompts-one-line-qa-2

evidence · `01M79C8NH4QJ98RDG6K99H76JS` · actor [claude-code-9b656f6b](../actor/claude-code-9b656f6b.md)

## result

PASS

## statement

`rhiza workflow describe build` (built at HEAD) shows the builder prompt `/rhiza-builder Build leaf {leaf} of change set {change} ({title}), attempt {attempt}: follow rhiza guide builder.` and the verifier one-liner; `workflow describe review` shows the same verifier line. .rhiza/workflows/docs-build.toml starts the builder prompt with the builtin line plus the humanize addendum; docs-review.toml's verifier prompt is byte-identical. go test ./internal/workflow/... ok; TestScripts/(build-evidence-commit|docs-commands-exist) ok

## Links

### Out

- verifies [build-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-prompts-one-line.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
