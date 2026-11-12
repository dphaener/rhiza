# verdict-downgraded

capability · `01M76XZGDTNAM0425G73YGZ3P0` · parent [review-critic](../review-critic.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run RunVerdict) assert a downgrade of a major to minor is journaled with both severities, and a downgrade to major of a major, to blocker, or with no severity exits as a usage error; go test -trimpath ./internal/run covers the journal entry; describe-golden is regenerated (go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden).

## statement

rhiza run verdict accepts a third verdict, downgraded, with a severity lower than the finding's (blocker to major, minor or nit; major to minor or nit), journaled with the finding's original severity; a downgrade to the same or a higher severity, or without a severity, is refused.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [verdict-downgraded-verified](../../../../../../evidence/verdict-downgraded-verified.md)
- verifies [verdict-downgraded-qa-1](../../../../../../evidence/verdict-downgraded-qa-1.md)
- verifies [pr-510-critic-downgrade](../../../../../../evidence/pr-510-critic-downgrade.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
