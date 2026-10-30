# review-code-check-in-rex

capability · `01M6505Z15CNR0WFDTKBVMHK70` · parent [run-program-status](../run-program-status.md)

## acceptance

A test runs a review whose code check fails through a fake Rex and finds the root record working with the code-check step, then error naming the failed run

## statement

While a review's code check runs, Rex's root record names the code-check step, and a failed code check ends the root record in error

## Links

### Out

- decided-by [run-reports-progress](../../../../../decision/run-reports-progress.md)

### In

- verifies [review-code-check-in-rex-verified](../../../../../evidence/review-code-check-in-rex-verified.md)
- verifies [review-code-check-in-rex-qa-1](../../../../../evidence/review-code-check-in-rex-qa-1.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
