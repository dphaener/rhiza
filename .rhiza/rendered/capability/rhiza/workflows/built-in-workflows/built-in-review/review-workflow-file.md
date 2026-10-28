# review-workflow-file

capability · `01M5Z0EPGNJBA661QHJ901DRYA` · parent [built-in-review](../built-in-review.md)

## acceptance

A testscript with fake agents runs `rhiza run start review` on a built change set: one check reports a major the critic confirms, a fixer and the leaf's verifier run, only that check re-runs clean, a fake QA check writes Evidence as its own actor and build-check still passes, `rhiza run approve` ends the run done and the change set stays accepted; describe-golden is regenerated.

## statement

Rhiza ships a built-in `review` workflow: the code check, then a checks step with compliance, quality, qa, scope and model checks and their critic on opus, and the sign-off, ending done without landing the change set.

## Links

### Out

- depends-on [review-check-worktree](review-check-worktree.md)
- depends-on [review-cycle-cap](review-cycle-cap.md)
- depends-on [review-code-check](review-code-check.md)
- depends-on [review-signoff](review-signoff.md)
- depends-on [review-summary](review-summary.md)
- depends-on [review-guide](review-guide.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [four-workflows](../../../../../decision/four-workflows.md)
- decided-by [review-feeds-model](../../../../../decision/review-feeds-model.md)
- decided-by [review-ruling-14](../../../../../decision/review-ruling-14.md)
- decided-by [review-ruling-23](../../../../../decision/review-ruling-23.md)
- decided-by [review-ruling-24](../../../../../decision/review-ruling-24.md)

### In

- verifies [review-workflow-file-verified](../../../../../evidence/review-workflow-file-verified.md)
- verifies [review-workflow-file-head-verified](../../../../../evidence/review-workflow-file-head-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
