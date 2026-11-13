# build-prompts-one-line-verified

evidence · `01M79AQW966Q58JG1W7FBM4TBW` · actor [claude-code-f52ff98e](../actor/claude-code-f52ff98e.md)

## result

PASS

## statement

Verified at b8abc7335 by claude-code-f52ff98e. Read the diff. build.toml's builder prompt is exactly the /rhiza-builder line and its verifier prompt is exactly the /rhiza-verifier line. review.toml's and docs-review.toml's verifier prompts equal that line byte for byte. docs-build.toml uses the builder line plus the doc-prose humanize addendum, which is unchanged from HEAD~1, and the verifier line. Ran go test -trimpath ./internal/workflow/...: ok. TestBuildBuiltin asserts both exact lines and that neither carries Checks:, session done, checklist, full suite or CI runs. TestReviewBuiltin passes. New TestBuildVerifierPromptsShared reads .rhiza/workflows/docs-build.toml and docs-review.toml and passes. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-evidence-commit|docs-commands-exist|describe-golden)': ok, with describe-golden unchanged by the commit. make self-check with the run's RHIZA_ROOT in a clean detached worktree at the commit failed only on render files for this verifier's own session actor (claude-code-f52ff98e), which is registered by the run but not committed; this is the known trap for run sessions. With the root set to the clean checkout, HEAD~1 fails the same actor render checks as this commit, so the leaf adds no failures. The builder reported make self-check exit 0.

## Links

### Out

- verifies [build-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-prompts-one-line.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
