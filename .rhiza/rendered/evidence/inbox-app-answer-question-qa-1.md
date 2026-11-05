# inbox-app-answer-question-qa-1

evidence · `01M69VVWCNCNTRCJ5GDM8XCBB2` · actor [claude-code-77fbcb0f](../actor/claude-code-77fbcb0f.md)

## result

PASS

## statement

At 302d9b36, drove the real inbox app (rhiza tui app inbox in a pty) on a scratch repo whose current change set was a draft and where store-q blocks an accepted change set: a opened the answer prompt on the blocker item, typing use sqlite then enter showed answer done and Nothing waits on you; rhiza show answer-store-q --change <accepted cs> showed the Decision quoting use sqlite, linked answers store-q, in the blocked change set rather than the current draft. In a copy, a, typing, then esc and r left the blocker listed and wrote no decision. The run-question session.message path needs a live Rex session and was covered only by the Go fake-runner test.

## Links

### Out

- verifies inbox-app-answer-question (retired)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
