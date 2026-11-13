# built-in-review

capability · `01M5Z0EFR7RVGGDDPQF4B9TQN4` · parent [built-in-workflows](../built-in-workflows.md)

## statement

The review workflow reviews a built change set: a code check first, then fresh check sessions for compliance (which also judges scope: under-built, scope creep and unaccounted files), quality, hands-on QA and the model, a critic that can't add findings settling failed checks, fixes routed per leaf at most twice, and the owner's sign-off with rhiza run approve.

## Links

### Out

- decided-by [review-workflow-shape](../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-3](../../../../decision/review-ruling-3.md)
- decided-by [review-ruling-4](../../../../decision/review-ruling-4.md)
- decided-by [review-ruling-5](../../../../decision/review-ruling-5.md)
- decided-by [review-ruling-6](../../../../decision/review-ruling-6.md)
- decided-by [review-ruling-7](../../../../decision/review-ruling-7.md)
- decided-by [review-ruling-27](../../../../decision/review-ruling-27.md)
- decided-by [review-role-topics-shape](../../../../decision/review-role-topics-shape.md)
- derived-from [review-workflow-plan](../../../../source/review-workflow-plan.md)

### In

- verifies [pr-45-the-review-workflow](../../../../evidence/pr-45-the-review-workflow.md)
- verifies [pr-537-review-role-topics](../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
