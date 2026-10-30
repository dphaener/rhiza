# review-questions-in-rex

capability · `01M6505ZVVRWAHQW8S8FGJ25MC` · parent [run-program-status](../run-program-status.md)

## acceptance

A test exhausts a review's two fix cycles through a fake Rex and finds the root record blocked, its message naming the change set, the review step and the question's text

## statement

When a review ends with open questions after its two fix cycles, Rex's root record is blocked and names the first open question with its change set and step

## Links

### Out

- decided-by [run-reports-progress](../../../../../decision/run-reports-progress.md)
- decided-by [run-ends-blocked-on-review-question](../../../../../decision/run-ends-blocked-on-review-question.md)

### In

- blocks [review-questions-child-record](../../../../../question/review-questions-child-record.md)
- verifies [review-questions-in-rex-verified](../../../../../evidence/review-questions-in-rex-verified.md)
- verifies [review-questions-in-rex-qa-1](../../../../../evidence/review-questions-in-rex-qa-1.md)
- verifies [review-questions-in-rex-verified-2](../../../../../evidence/review-questions-in-rex-verified-2.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
