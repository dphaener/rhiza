# conduct-no-wait-guide

capability · `01M7QE205YR8ZF9RC2WQ469440` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a test asserting conduct-roles.md documents --no-wait and the wake-up message and conduct.md tells the conductor to end its turn after a --no-wait spawn; (failure case) the test fails on the current text, which says only that the spawn waits for the child. go test -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestScripts/guide-install' passes.

## files

internal/guide/text/conduct.md, internal/guide/text/conduct-roles.md, internal/guide/guide_test.go, docs/primitives.md, docs/bring-your-own-agent.md

## must_not

Change any role other than the conductor's instructions; drop the rule that the conductor never sleeps or polls; tell a conductor to poll run status for its children

## statement

The conduct guides (conduct, conduct-roles) and the run spawn docs tell a conductor to spawn children with --no-wait, start independent children together, end its turn, and act on the wake-up session message each child's end sends (reading rhiza run status --json for the outcomes), and still describe the default waiting spawn for other callers.

## Links

### Out

- depends-on [spawn-wake-conductor](spawn-wake-conductor.md)
- decided-by [spawn-no-wait-wake](../../../../decision/spawn-no-wait-wake.md)

### In

- verifies [conduct-no-wait-guide-verified](../../../../evidence/conduct-no-wait-guide-verified.md)
- verifies [conduct-no-wait-guide-qa-1](../../../../evidence/conduct-no-wait-guide-qa-1.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
