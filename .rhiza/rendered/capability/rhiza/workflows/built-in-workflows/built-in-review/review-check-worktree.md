# review-check-worktree

capability · `01M5Z0EGRD1KYTEXXSWDEZNJCE` · parent [built-in-review](../built-in-review.md)

## acceptance

A test finds a check session's working directory a detached worktree at the branch head, gone after the session ends whatever its outcome, a finding or Evidence it wrote in the run worktree's .rhiza/log or journal, and uncommitted edits in the run worktree untouched.

## statement

Each check and critic session works in a throwaway detached worktree at the run branch's head commit, removed once its session ends, while its rhiza writes go into the run worktree's log.

## Links

### Out

- depends-on [review-checks-step](review-checks-step.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-ruling-13](../../../../../decision/review-ruling-13.md)
- decided-by [review-ruling-14](../../../../../decision/review-ruling-14.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- verifies [review-check-worktree-verified](../../../../../evidence/review-check-worktree-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
