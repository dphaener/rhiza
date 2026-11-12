# severity-contract-qa-1

evidence · `01M78XSAHD0VY4RJW8HNC2CEAZ` · actor [claude-code-9af4801e](../actor/claude-code-9af4801e.md)

## result

PASS

## statement

QA at e313e4634 with the built binary: I extracted '## Severities' from rhiza guide review and rhiza guide plan-review, and cmp found them byte-identical (blocker, major, minor and nit, a per-check mapping and a plan-review mapping). The block names no go or make command. Both Critic sections say confirm, downgrade or drop and never add a finding or raise a severity, and both show the downgraded verdict's severity field.

## Links

### Out

- verifies [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
