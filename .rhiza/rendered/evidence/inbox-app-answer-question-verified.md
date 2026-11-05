# inbox-app-answer-question-verified

evidence · `01M69VMQ8PHFVCBFQ84D7303AH` · actor [claude-code-f72c0518](../actor/claude-code-f72c0518.md)

## result

PASS

## statement

At 6f7c433d ran go test ./internal/tui/inboxapp -count=1 -v: TestAnswerBlockerCallsQuestionAnswer (fake runner) records exactly one question.answer call with args which-db|C9|use sqlite (question handle, change set, typed words) for a blocker item; TestAnswerRunQuestionStillMessagesSession records one session.message SESS1|sqlite for a run question; TestEscOnBlockerAnswerSendsNothing records no calls after esc. All pass, as do go test ./..., go vet and fmt-check. The diff adds a Blocking questions group, enables a on blocker items, and routes blocker answers to question.answer. Note: question.answer is not registered yet (leaf question-answer-op, unbuilt); the app passes change positionally, so that op must accept <q> <cs> <words> as argv or this call path errors in the real app.

## Links

### Out

- verifies inbox-app-answer-question (retired)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
