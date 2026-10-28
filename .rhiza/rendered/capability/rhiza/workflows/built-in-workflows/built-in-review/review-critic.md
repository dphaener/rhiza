# review-critic

capability · `01M5Z0EHF5WMV1Y0X96EJ0NCPZ` · parent [built-in-review](../built-in-review.md)

## acceptance

An engine test with fake sessions finds no critic for a clean check or one with only minors and nits and both passed, a critic for a check with a major that drops it and the check passed, and a critic confirming a blocker and the check failed; a verdict adding a finding is refused.

## statement

After a cycle the engine starts a fresh critic session, in parallel, only for each check whose latest round has a blocker or major finding, and the critic can only confirm or drop those findings; a check fails only when a confirmed blocker or major survives, and minor and nit findings never call a critic nor fail a check.

## Links

### Out

- depends-on [review-checks-step](review-checks-step.md)
- depends-on [review-findings](review-findings.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [critics-on-serious-findings](../../../../../decision/critics-on-serious-findings.md)
- decided-by [review-ruling-10](../../../../../decision/review-ruling-10.md)
- decided-by [review-ruling-15](../../../../../decision/review-ruling-15.md)
- decided-by [review-ruling-16](../../../../../decision/review-ruling-16.md)

### In

- depends-on [review-fix-leaf](review-fix-leaf.md)
- depends-on [review-fix-change-set](review-fix-change-set.md)
- verifies [review-critic-verified](../../../../../evidence/review-critic-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
