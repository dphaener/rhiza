# review-signoff

capability · `01M5Z0EMM6HBY40F6M5QVK0TX0` · parent [built-in-review](../built-in-review.md)

## acceptance

Testscripts find a review run waiting at sign-off until `rhiza run approve` and then done with an approved entry naming the actor; `run request-changes` with one finding starts a fixer for the named leaf after the cap was used and returns to the wait; outside a waiting review run both commands exit 2; an autonomous run passes sign-off with a journal entry of rule autonomous; describe-golden is regenerated.

## statement

The review's last step waits for the owner: `rhiza run approve` journals the owner's sign-off and the run ends done, and `rhiza run request-changes <file>` journals the owner's findings as already confirmed and starts one fix cycle outside the cap, after which the failed checks re-run and the run waits for the owner again; an autonomous run skips the wait and journals the skip with rule autonomous.

## Links

### Out

- depends-on [review-fix-leaf](review-fix-leaf.md)
- depends-on [review-cycle-cap](review-cycle-cap.md)
- decided-by [two-owner-gates](../../../../../decision/two-owner-gates.md)
- decided-by [autonomous-mode](../../../../../decision/autonomous-mode.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-8](../../../../../decision/review-ruling-8.md)
- decided-by [review-ruling-22](../../../../../decision/review-ruling-22.md)
- decided-by [review-ruling-23](../../../../../decision/review-ruling-23.md)
- decided-by [review-ruling-26](../../../../../decision/review-ruling-26.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- depends-on [review-summary](review-summary.md)
- verifies [review-signoff-verified](../../../../../evidence/review-signoff-verified.md)
- verifies [review-signoff-lock-verified](../../../../../evidence/review-signoff-lock-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
