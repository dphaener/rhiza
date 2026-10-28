# review-ruling-8

decision · `01M5ZQFPCGMFAN7TERGKYGPE7D`

## options

```json
[
  "New signoff step kind (Recommended)",
  "Tail of the checks step",
  "Wait step + command"
]
```

## picked

New signoff step kind (Recommended)

## question

How does sign-off sit in the workflow?

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

How does sign-off sit in the workflow? — New signoff step kind (Recommended). A signoff step has outcomes pass and changes, its on.changes must name a checks step, and it waits by stopping the run until run approve or run request-changes resumes it.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [review-signoff](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-signoff.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
