# downgrade-ends-fix-verified

evidence · `01M78XDKQ2V058G8NZ1N468KDS` · actor [claude-code-7883df3c](../actor/claude-code-7883df3c.md)

## result

PASS

## statement

At 96c2f0df4 ran go test -trimpath -count=1 ./internal/run -run 'Checks|Critic' (ok): TestCriticDowngradeEndsFix has the critic downgrade alpha's only major to minor and asserts the review passes, only checks and the critic start, r.fixSpecs and r.leafSpecs are empty, and the outcome detail is 'cycle 1: critic downgraded f1 (major to minor)'; TestCriticConfirmedMajorStartsFix confirms one major and downgrades another, and asserts the check fails and exactly one change-set fixer starts with only 'alpha f1 major'. Ran go test -trimpath -count=1 ./internal/api -run 'RunStatus|ReviewSummary' (ok): TestRunStatusReviewSummaryDowngraded asserts the summary lists the finding under 'Downgraded by the critic' as 'downgraded from major to minor', not under Confirmed, and that the run status finding keeps severity major with downgraded_to minor. Read the diff (checks.go checkVerdict, run_review.go, run_summary.go): it matches. Both attempt-1 findings (fixer never checked; no confirmed-major fixer case) are fixed.

## Links

### Out

- verifies [downgrade-ends-fix](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/downgrade-ends-fix.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
