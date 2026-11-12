# severity-contract

capability · `01M76XZGFGB9C8J6GA5NF3RBQE` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'Review|Plan' asserts review.md and plan-review.md contain the same severity block byte for byte, that both critic sections name confirmed, downgraded and dropped, and that the block names no Go, make or Rhiza-repo command.

## statement

The review and plan-review guide topics hold one shared severity block, byte for byte: blocker (behaviour broken, a build or test fails, or a leaf is missing), major (an acceptance point unmet or a real defect), minor (style, naming, docs or a non-blocking follow-up) and nit, with a mapping per check; and the critic sections say confirm, downgrade or drop, never add or raise.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [severity-contract-verified](../../../../../../evidence/severity-contract-verified.md)
- verifies [severity-contract-qa-1](../../../../../../evidence/severity-contract-qa-1.md)
- verifies [severity-contract-verified-2](../../../../../../evidence/severity-contract-verified-2.md)
- verifies [severity-contract-verified-5](../../../../../../evidence/severity-contract-verified-5.md)
- verifies [pr-510-critic-downgrade](../../../../../../evidence/pr-510-critic-downgrade.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
