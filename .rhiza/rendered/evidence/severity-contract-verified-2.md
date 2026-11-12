# severity-contract-verified-2

evidence · `01M78YPWCMVS2BPJ7CD5G75JDB` · actor [claude-code-76b815e0](../actor/claude-code-76b815e0.md)

## result

PASS

## statement

At 49cf7e8ad ran go test -trimpath ./internal/guide -run 'Review|Plan' -v: all pass, including TestReviewSeverityBlockShared, which asserts the review and plan-review '## Severities' blocks are byte-identical, that the block names no go/Go/make/rhiza check/render/ship/testscript, and that both Critic sections contain confirmed, downgraded, dropped and 'Never add a finding'. Read the block: blocker/major/minor/nit plus a per-check mapping (plan review: missing-acceptance, contradiction, outside-intent major, the rest minor). The attempt-1 finding (plan-review downgrade always refused) is fixed: run.KindSeverity/Finding.EffectiveSeverity give a kind its mapped severity, runRunVerdict, review-check and the gate use it (Round.Blocking); go test -trimpath ./internal/run/... and ./internal/api -run 'PlanReviewDowngradeByKind|RunVerdict|ReviewFindingsNameLinks|ReviewCheck|Gate' pass.

## Links

### Out

- verifies [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
