# review-code-check

capability · `01M5Z0EKX3DVJHR6SVG5QX87J1` · parent [built-in-review](../built-in-review.md)

## acceptance

Testscripts find code-check passing on a change set whose leaves have verifier Evidence and a clean model, and failing, naming the cause, on one with a leaf lacking Evidence and on one where rhiza check reports a problem; a workflow override adding a command step after it runs that command before the checks step.

## statement

`rhiza run code-check` runs `rhiza check` and `rhiza run build-check` on the run's change set and journals the result as the review's code check, and the built-in review runs it as a command step before any agent check, which a repo's workflow override can replace or extend with its own test command.

## Links

### Out

- depends-on [review-checks-step](review-checks-step.md)
- decided-by [code-checks-first](../../../../../decision/code-checks-first.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-1](../../../../../decision/review-ruling-1.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- depends-on [code-check-skips-reference](code-check-skips-reference.md)
- blocks [review-code-check-fails](../../../../../question/review-code-check-fails.md)
- verifies [review-code-check-verified](../../../../../evidence/review-code-check-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
