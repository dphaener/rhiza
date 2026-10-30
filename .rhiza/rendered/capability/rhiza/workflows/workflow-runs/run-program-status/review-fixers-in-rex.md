# review-fixers-in-rex

capability · `01M6505Y4EB35E286KJRDR8JBV` · parent [run-program-status](../run-program-status.md)

## acceptance

A test runs a review with one confirmed finding through a fake Rex and finds, in order, a child record for the leaf fixer, then the leaf's verifier, then a child record for the change-set fixer, each cleared once its session ends, with the root record naming the review step

## statement

While a review run goes, Rex shows each live fixer and verifier session of the review as a child record: a leaf fixer or its verifier names its leaf, and the change-set fixer names the review

## Links

### Out

- decided-by [run-reports-progress](../../../../../decision/run-reports-progress.md)

### In

- verifies [review-fixers-in-rex-verified](../../../../../evidence/review-fixers-in-rex-verified.md)
- verifies [review-fixers-in-rex-qa-1](../../../../../evidence/review-fixers-in-rex-qa-1.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
