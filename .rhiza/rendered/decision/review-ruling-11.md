# review-ruling-11

decision · `01M5ZQFPVRDXEF6N3X7R36Z2SK`

## options

```json
[
  "[[step.X.check]] array (Recommended)",
  "[step.X.checks.<name>] map"
]
```

## picked

[[step.X.check]] array (Recommended)

## question

How are checks declared in TOML?

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

How are checks declared in TOML? — [[step.X.check]] array (Recommended). Checks are an array of tables, [[step.X.check]] with name = "...", run in file order.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [review-checks-step](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
