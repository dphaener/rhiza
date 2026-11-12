# pr-510-critic-downgrade

evidence · `01M790C2PDGKJ612FH6N6SH8QQ` · actor [dphaener](../actor/dphaener.md)

## head

d4617bde3fc46dccc9114c890c9926c6ff5fcab8

## pr

https://github.com/dphaener/rhiza/pull/510

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #510 (https://github.com/dphaener/rhiza/pull/510) passed CI at d4617bd

## status

verified

## Links

### Out

- verifies [verdict-downgraded](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/verdict-downgraded.md)
- verifies [downgrade-ends-fix](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/downgrade-ends-fix.md)
- verifies [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)
- verifies [plan-review-rounds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-rounds.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
