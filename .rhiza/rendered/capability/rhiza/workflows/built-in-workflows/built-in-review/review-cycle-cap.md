# review-cycle-cap

capability · `01M5Z0EK65NC55WAS23EJDV773` · parent [built-in-review](../built-in-review.md)

## acceptance

An engine test whose check fails three times finds two fix cycles, a question blocking the named leaf, no sign-off step entered and the run failed; after an answering Decision a resume runs one more cycle with the Decision in the fixer's prompt; in an autonomous run the arbiter is started on the question and a further failure ends the run failed.

## statement

A review runs at most 2 fix cycles; after that each surviving confirmed finding becomes a question, linked blocks to the leaf it names, and the run ends failed before sign-off; a resume once they are answered gives one fresh cycle with the answering Decisions in the fixers' prompts, and in an autonomous run the arbiter answers at once and the review gets that one fresh cycle, failing the run if a check fails again.

## Links

### Out

- depends-on [review-fix-leaf](review-fix-leaf.md)
- depends-on [review-fix-change-set](review-fix-change-set.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [arbiter](../../../../../decision/arbiter.md)
- decided-by [review-ruling-20](../../../../../decision/review-ruling-20.md)
- decided-by [review-ruling-21](../../../../../decision/review-ruling-21.md)
- decided-by [review-ruling-22](../../../../../decision/review-ruling-22.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- depends-on [review-signoff](review-signoff.md)
- verifies [review-cycle-cap-verified](../../../../../evidence/review-cycle-cap-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
