# ask-status-is-question

capability · `01M6505W1FFMVWMMDTMFDKA9CN`

## acceptance

cmd/rhiza/testdata/script/run-blocked.txtar passes: a run's AskUserQuestion reports state=blocked:kind=question with the question as msg

## statement

Rex program status shows a session's AskUserQuestion as blocked with kind question and the question as its message, never as a permission

## Links

### Out

- decided-by [rex-notice-by-program-status](../decision/rex-notice-by-program-status.md)
- derived-from [questions-not-permissions-plan-run](../source/questions-not-permissions-plan-run.md)

### In

- blocks [rex-ask-wording](../question/rex-ask-wording.md)
- verifies [ask-status-is-question-verified](../evidence/ask-status-is-question-verified.md)
- verifies [ask-status-is-question-qa-1](../evidence/ask-status-is-question-qa-1.md)
- verifies [ask-status-is-question-verified-2](../evidence/ask-status-is-question-verified-2.md)
- verifies [pr-115-questions-not-permissions](../evidence/pr-115-questions-not-permissions.md)

## History

- `01M61GW8H09EPJ1YCR7P2X3E8J` Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.
