# review-fix-change-set

capability · `01M5Z0EJRDQ0F1VZBYFM5K6DXX` · parent [built-in-review](../built-in-review.md)

## acceptance

An engine test finds two confirmed findings naming no leaf handed to a single fixer session, a new commit with a `Review:` trailer naming the run, and the failed checks re-run after it; a fixer ending without such a commit fails the cycle.

## statement

Confirmed blocker or major findings that name no leaf go to one fixer session for the change set, whose fix commit carries a `Review:` trailer naming the run, before the failed checks re-run.

## Links

### Out

- depends-on [review-critic](review-critic.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-18](../../../../../decision/review-ruling-18.md)
- decided-by [review-ruling-19](../../../../../decision/review-ruling-19.md)

### In

- depends-on [review-cycle-cap](review-cycle-cap.md)
- verifies [review-fix-change-set-verified](../../../../../evidence/review-fix-change-set-verified.md)
- verifies [review-fix-change-set-head-verified](../../../../../evidence/review-fix-change-set-head-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
