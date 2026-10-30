# review-questions-in-rex-qa-1

evidence · `01M66DH44JX9Q3B3NKAZM2KC76` · actor [claude-code-31db92cd](../actor/claude-code-31db92cd.md)

## result

PASS

## statement

Built ./cmd/rhiza at 5dd508d and ran the builtin review inside a real Rex block on a scratch repo where a fake quality check confirmed a major on leaf alpha in every cycle. Rex showed alpha's fixer and verifier at attempts 2 and 3, and quality at cycles 2 and 3. Then the root record was blocked: 'Review me: step review waits on you: Review check quality still fails after 3 cycles: finding f1 (major at alpha.txt:1) on leaf alpha: still too short. ...'. That message names the change set, the review step and the question's text. rhiza run status agreed: an unanswered question after cycle 3 on f1, blocking leaf alpha.

## Links

### Out

- verifies [review-questions-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-questions-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
