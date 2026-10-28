# review-checks-step

capability · `01M5Z0EG65NMPZAKN41R95EQYB` · parent [built-in-review](../built-in-review.md)

## acceptance

Parse tests load a checks step and reject one with no checks, a duplicate check name or an unknown sub-table key; an engine test with a fake runner over three checks finds all three sessions started before any ends, and in a second cycle after one check failed finds only that check's session started; `rhiza run status` shows each check's cycles and outcome.

## statement

Workflow files take a `checks` step kind naming its checks (each with prompt, agent, permission_mode, models) and `critic` and `fixer` sub-tables, and the engine starts every check as a fresh session in parallel, each in its own pane of the change set's tab told only the change set ID and its check name, journaling each check's session and outcome; a later cycle re-runs only the checks that failed.

## Links

### Out

- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [re-review-failures-only](../../../../../decision/re-review-failures-only.md)
- decided-by [full-agent-sessions](../../../../../decision/full-agent-sessions.md)
- decided-by [review-ruling-9](../../../../../decision/review-ruling-9.md)
- decided-by [review-ruling-11](../../../../../decision/review-ruling-11.md)
- decided-by [review-ruling-12](../../../../../decision/review-ruling-12.md)
- decided-by [review-ruling-24](../../../../../decision/review-ruling-24.md)

### In

- depends-on [review-check-worktree](review-check-worktree.md)
- depends-on [review-critic](review-critic.md)
- depends-on [review-code-check](review-code-check.md)
- verifies [review-checks-step-verified](../../../../../evidence/review-checks-step-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
