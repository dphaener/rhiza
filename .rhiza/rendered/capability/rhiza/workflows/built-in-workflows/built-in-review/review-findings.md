# review-findings

capability · `01M5Z0EH5ZG74DZQY63G9B78WR` · parent [built-in-review](../built-in-review.md)

## acceptance

A testscript in a review run records findings for two checks and a verdict on one: the journal holds them per check and `run status --json` shows them; a finding with an unknown severity, an unknown check, a leaf not in the run's change set, or a verdict naming a finding its check's latest round does not hold exits 2 with nothing written; the plan-review testscripts still pass; describe-golden is regenerated.

## statement

In a review run `rhiza run findings` and `rhiza run verdict` take the check they belong to, and each review finding has a severity (blocker, major, minor, nit), the leaf it is about when there is one, a place (file:line or node) and a reason, while plan review findings keep their five kinds; `rhiza run status` shows each check's rounds with findings and verdicts.

## Links

### Out

- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-12](../../../../../decision/review-ruling-12.md)

### In

- depends-on [review-critic](review-critic.md)
- depends-on [review-summary](review-summary.md)
- depends-on [review-guide](review-guide.md)
- verifies [review-findings-verified](../../../../../evidence/review-findings-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
