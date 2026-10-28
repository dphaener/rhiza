# review-fix-leaf

capability · `01M5Z0EJ62SZXJR5NSAA9X48ZV` · parent [built-in-review](../built-in-review.md)

## acceptance

An engine test against a real git repo finds, for a confirmed finding naming leaf A, a fixer whose prompt names A and the finding, a new commit with `Leaf: A`, A's verifier started on it and its Evidence written by the verifier's actor, then only the failed check re-run; a fixer ending with no commit naming A fails the cycle.

## statement

Each leaf with confirmed blocker or major findings goes to a fresh fixer session that reads them from the journal and adds a fix commit with that leaf's `Leaf:` trailer, after which that leaf's verifier runs again as in the build and writes fresh Evidence, and only then do the failed checks re-run.

## Links

### Out

- depends-on [review-critic](review-critic.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- decided-by [review-ruling-17](../../../../../decision/review-ruling-17.md)
- decided-by [review-ruling-18](../../../../../decision/review-ruling-18.md)
- decided-by [review-ruling-20](../../../../../decision/review-ruling-20.md)

### In

- depends-on [review-cycle-cap](review-cycle-cap.md)
- depends-on [review-signoff](review-signoff.md)
- verifies [review-fix-leaf-verified](../../../../../evidence/review-fix-leaf-verified.md)
- verifies [review-fix-leaf-head-verified](../../../../../evidence/review-fix-leaf-head-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
