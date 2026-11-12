# severity-contract-verified-5

evidence · `01M78ZG2HV7RMBAP5CSTWG56G1` · actor [claude-code-63175055](../actor/claude-code-63175055.md)

## result

PASS

## statement

Verifier, attempt 5, at d05bba3b8. Ran the acceptance: go test -trimpath -count=1 ./internal/guide -run 'Review|Plan' passes; TestReviewSeverityBlockShared asserts review.md and plan-review.md hold the same Severities block byte for byte, that it names no Go, make, rhiza check/render/ship or testscript command, and that both Critic sections name confirmed, downgraded and dropped and say never add a finding. Read the diff: the block lists blocker, major, minor and nit with a mapping per check. The earlier findings are fixed: plan-review findings take their kind's severity (run.KindSeverity), so a plan critic's downgrade works, and review-check and the gate count only Round.Blocking. plan-review-rounds is amended to match. Also passed: go test -trimpath ./internal/api -run 'TestPlanReviewDowngradeByKind|Gate|ReviewCheck|Verdict|Downgrad', ./internal/run/..., the cmd/rhiza run-review*, review-*, plan-review-visit-round and autonomous-review scripts, and make self-check in a clean detached checkout with the run's RHIZA_* env unset (exit 0). In this worktree it fails only on this session's own unrendered actor record.

## Links

### Out

- verifies [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
