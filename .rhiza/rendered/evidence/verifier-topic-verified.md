# verifier-topic-verified

evidence · `01M799QRXRQ6YK5DCJ80RZQVC3` · actor [claude-code-457b755e](../actor/claude-code-457b755e.md)

## result

PASS

## statement

At efc2f11f8 (detached) re-ran every Checks: command: go test -trimpath -count=1 ./internal/guide/... exit 0, ./internal/api -run Guide exit 0, ./cmd/rhiza -run TestScripts/guide exit 0, rhiza ship generate exit 0 (0 files changed), rhiza render exit 0 (no tracked changes). Acceptance go test -trimpath ./internal/guide/... passes. Read the diff: internal/guide/verifier_test.go TestVerifierText checks the six headings in order, the session-done, test-scope, command-budget and never-list shared blocks expanded, detached checkout, Never edit code, fakes only, the Checks: and acceptance re-run with the missing-block/different-exit-code finding naming the command (TestBuildVerifierRerunsChecks removed from build_test.go), the scope and no-stubs points, the Evidence apply example, the rhiza run findings heredoc and the finding-names-file/command/acceptance-point rule, plus a negative case banning builder commit/trailer instructions. rhiza guide verifier prints the expanded text; guide ls count and List order updated. No stubs or changes outside the leaf.

## Links

### Out

- verifies [verifier-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-topic.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
