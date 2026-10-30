# run-ends-blocked-on-review-question

decision · `01M66DR7ZX5M1M4HMYZXVVV7VV`

## statement

A run reports its progress to Rex through OSC 7501 program status: the run is the root record, with its current step and progress as leaves done out of all; each live leaf is a child record with a hierarchical ID, the builder or verifier working or blocked; and the run ends done when it succeeds, or error when it fails, except that a failed run whose ending step leaves a review question open ends blocked, naming that question.

## wording

run-reports-progress's statement, amended by the review-questions-in-rex fixer for the blocked ending the autonomous plan run for change set 01M602TG planned; no owner words

## Links

### Out

- supersedes [run-reports-progress](run-reports-progress.md)
- derived-from [review-status-plan-autonomous](../source/review-status-plan-autonomous.md)

### In

- decided-by [run-program-status](../capability/rhiza/workflows/workflow-runs/run-program-status.md)
- decided-by [review-questions-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-questions-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
