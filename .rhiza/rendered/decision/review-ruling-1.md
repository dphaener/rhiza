# review-ruling-1

decision · `01M5ZQFN92JF2814QMT3YJGTMX`

## options

```json
[
  "Run ends failed at once (Recommended)",
  "Blocker to the fix cycle"
]
```

## picked

Run ends failed at once (Recommended)

## question

Question review-code-check-fails: when the review's code check (rhiza check + build-check) fails, should the run end failed at once?

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

Question review-code-check-fails: when the review's code check (rhiza check + build-check) fails, should the run end failed at once? — Run ends failed at once (Recommended). The built-in review's code-check step has no on.fail, so a failed code check ends the run failed, and a repo override may route it elsewhere.

## wording

picked by the orchestrator acting as arbiter

## Links

### Out

- answers [review-code-check-fails](../question/review-code-check-fails.md)

### In

- decided-by [review-code-check](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-code-check.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
