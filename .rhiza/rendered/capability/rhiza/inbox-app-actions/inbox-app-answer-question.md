# inbox-app-answer-question

capability · `01M686TSHQQ5XJ0YAZ3Z538821` · parent [inbox-app-actions](../inbox-app-actions.md)

## acceptance

A Go test in internal/tui/inboxapp with a fake runner asserts one question.answer call with the question, change set and typed words for a blocker item, that an answer on a run question still calls session.message, and that esc sends nothing

## statement

In the inbox app, typing an answer on a blocker item answers it in place through question.answer, and the item is gone on reload

## Links

### Out

- decided-by [inbox-blockers-per-change-set](../../../decision/inbox-blockers-per-change-set.md)

### In

- verifies [inbox-app-answer-question-verified](../../../evidence/inbox-app-answer-question-verified.md)
- verifies [inbox-app-answer-question-qa-1](../../../evidence/inbox-app-answer-question-qa-1.md)
- verifies [pr-143-inbox-shows-blockers](../../../evidence/pr-143-inbox-shows-blockers.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
