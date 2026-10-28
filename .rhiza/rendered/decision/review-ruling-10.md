# review-ruling-10

decision · `01M5ZQFPPN4SK89F9MW1MNEJ59`

## options

```json
[
  "Fails, re-runs next cycle (Recommended)",
  "Passes as clean"
]
```

## picked

Fails, re-runs next cycle (Recommended)

## question

A check session ending without any findings round?

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

A check session ending without any findings round? — Fails, re-runs next cycle (Recommended). A check session that records no findings round fails its check, gets no critic and no fixer, and re-runs in the next cycle.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [review-critic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
