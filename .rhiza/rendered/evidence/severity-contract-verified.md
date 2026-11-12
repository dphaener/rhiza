# severity-contract-verified

evidence · `01M78X93F004AR792NRF3DJ512` · actor [claude-code-5524f454](../actor/claude-code-5524f454.md)

## result

PASS

## statement

At 9e405b92b, go test -trimpath ./internal/guide -run 'Review|Plan' -v passed all 13 tests, including the new TestReviewSeverityBlockShared, which extracts each topic's '## Severities' section and asserts review.md and plan-review.md match byte for byte, that both Critic sections contain confirmed, downgraded, dropped and 'Never add a finding', and that the block contains none of 'go ', 'make ', 'rhiza check/render/ship' or 'testscript'. Reading the diff: the shared block lists blocker (behaviour broken, build or test fails, leaf missing), major (acceptance point unmet, real defect), minor (style, naming, docs, non-blocking follow-up) and nit, with a mapping per check and for plan review; both critic sections say confirm, downgrade or drop and never add a finding or raise a severity. go test -trimpath ./internal/guide/... also passed. Side changes: TopicCap rose from 800 to 1000, and review.md's opener no longer says review verifiers get RHIZA_LEAF, although internal/run/checks.go still starts review verifiers. That is a minor docs loss and does not affect acceptance.

## Links

### Out

- verifies [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
