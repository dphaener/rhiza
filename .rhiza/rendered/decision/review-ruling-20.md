# review-ruling-20

decision · `01M602SXKWD2MFKKKEZ2W5HMJB`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

The fix phase belongs to cycle N+1: it runs after that cycle's check-cycle entry and before its check sessions, and a fixer's RHIZA_CYCLE is N+1. A failed fixer or verifier fails every check of that cycle with a fix-failed detail, so a failed fix uses up the cycle.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### In

- decided-by [review-fix-leaf](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-leaf.md)
- decided-by [review-cycle-cap](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-cycle-cap.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
