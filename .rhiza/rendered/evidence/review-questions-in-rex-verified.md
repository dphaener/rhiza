# review-questions-in-rex-verified

evidence · `01M66D9J248F171EYD3TAG7VP8` · actor [claude-code-7531ebc9](../actor/claude-code-7531ebc9.md)

## result

PASS

## statement

At 3b97350, go test ./internal/api -run TestProgramStatus -count=1 -v passes. TestProgramStatusReviewQuestions has check alpha record finding f1 in every cycle, so the review goes through fixer attempts 2 and 3 (its two fix cycles) and check cycle 3. Rex reports go to a buffer standing in for Rex. The last report is the root record (no id) with state=blocked and msg "Status runs: step review waits on you: Review check alpha still fails after 3 cycles: finding f1 (major): bad f1. ...". That names the change set by its title (rexRoot falls back to the change ID), the review step and the question's text. In the diff, deriveStatus maps a failed run whose ending step's detail lists unanswered questions to StatusBlocked, and the journal now keeps the question text in Detail. make check exits 0.

## Links

### Out

- verifies [review-questions-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-questions-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
