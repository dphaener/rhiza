# build-prompts-one-line-qa-1

evidence · `01M79B5586744PZAWTNQPBTV7C` · actor [claude-code-fb58129d](../actor/claude-code-fb58129d.md)

## result

PASS

## statement

Built rhiza; `rhiza workflow describe build` and `docs-build` show the one-line /rhiza-builder and /rhiza-verifier prompts, docs-build keeping the humanize addendum after the builder line; review.toml and docs-review.toml verifier prompts are byte-identical to build.toml's; `rhiza workflow check` reports build, review, docs-build, docs-review valid; go test ./internal/workflow/... and TestScripts/(build-evidence-commit|docs-commands-exist|describe-golden) ok; make self-check passes from a clean detached checkout.

## Links

### Out

- verifies [build-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-prompts-one-line.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
