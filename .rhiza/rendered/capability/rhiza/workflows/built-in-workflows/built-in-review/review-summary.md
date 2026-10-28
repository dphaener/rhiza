# review-summary

capability · `01M5Z0ENBN0TY0Z6FXG03NQB0P` · parent [built-in-review](../built-in-review.md)

## acceptance

A testscript renders the summary of a journaled review with a passed check, a check that failed and was fixed in one cycle, and two minors, and compares it with a golden file; `run status` at the sign-off wait names `rhiza run summary`.

## statement

`rhiza run summary` renders a review run's verdicts as markdown from the journal: each check's outcome, its confirmed findings, the minors and nits, the fix cycles and the sign-off, and `run status` names it while the run waits at sign-off.

## Links

### Out

- depends-on [review-findings](review-findings.md)
- depends-on [review-signoff](review-signoff.md)
- decided-by [one-pr-summary](../../../../../decision/one-pr-summary.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-16](../../../../../decision/review-ruling-16.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- verifies [review-summary-verified](../../../../../evidence/review-summary-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
