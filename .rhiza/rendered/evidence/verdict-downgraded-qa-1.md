# verdict-downgraded-qa-1

evidence · `01M78XS9ZZGF88PBM20B9FX2GM` · actor [claude-code-9af4801e](../actor/claude-code-9af4801e.md)

## result

PASS

## statement

QA at e313e4634 with the built binary (go build -trimpath ./cmd/rhiza), in a scratch Rhiza root running a real rhiza run start --foreground review workflow (fake rex and claude only stand in for agent sessions). The beta critic first sent four bad verdicts: major to major, blocker to blocker, downgraded with no severity, and a downgrade of a minor. Each exited 2 with a usage message (e.g. 'finding f2 is major: downgrade it to a lower severity than major, not major'). It then sent blocker to major and major to nit, which were accepted. The run journal holds {verdict:downgraded, severity:minor, original_severity:major} and the blocker to major and major to nit entries, each with original_severity.

## Links

### Out

- verifies [verdict-downgraded](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/verdict-downgraded.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
