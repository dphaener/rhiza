# downgrade-ends-fix

capability · `01M76XZGEMFYQN2SKCNBTHW3AT` · parent [review-critic](../review-critic.md)

## acceptance

Go tests in internal/run (go test -trimpath ./internal/run -run 'Checks|Critic') run a checks step where a check's only major is downgraded to minor and assert the check passes with no fixer session started, while a confirmed major still starts the fixer; go test -trimpath ./internal/api -run 'RunStatus|ReviewSummary' asserts both severities are shown.

## statement

A downgraded finding counts at its lower severity: a check whose surviving findings are all minor or nit after the critic passes its cycle without a fixer, and rhiza run status and the review summary show each downgraded finding with its original and lower severity.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [downgrade-ends-fix-verified](../../../../../../evidence/downgrade-ends-fix-verified.md)
- verifies [downgrade-ends-fix-qa-1](../../../../../../evidence/downgrade-ends-fix-qa-1.md)
- verifies [pr-510-critic-downgrade](../../../../../../evidence/pr-510-critic-downgrade.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
