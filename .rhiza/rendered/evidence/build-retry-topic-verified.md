# build-retry-topic-verified

evidence · `01M79A0ZAEK8K9ZK6KS0N2GWHB` · actor [claude-code-82bfd8da](../actor/claude-code-82bfd8da.md)

## result

PASS

## statement

Verifier at eba9d2913 (detached checkout): re-ran every Checks: command; go test -trimpath -count=1 ./internal/guide/... exit 0 (TestBuildRetryText PASS), ./internal/api -run Guide exit 0, ./cmd/rhiza -run TestScripts/guide exit 0, rhiza ship generate exit 0 (0 changed), rhiza render exit 0 (only this verifier's own actor rendered). Read the diff: internal/guide/text/build-retry.md registered as topic build-retry; the test asserts inspect the working tree + earlier attempts' commits, `rhiza run status --json`, Address every finding exactly as written, one fix commit with the same `Leaf:` trailer, Never rewrite history; asserts tokens <= TopicCap/3 (192 of 333) and that none of the builder's shared blocks (test-scope, unstable-ci, command-budget, outside-work, never-list, session-done) are repeated.

## Links

### Out

- verifies [build-retry-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-retry-topic.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
