# inbox-app-answer-verified-2

evidence · `01M6ND1CR45ENQSSNDS6CJ0C65` · actor [claude-code-d362bf4b](../actor/claude-code-d362bf4b.md)

## result

PASS

## statement

Verifier claude-code-d362bf4b at 8c79369c: go test -trimpath -count=1 ./internal/tui/inboxapp/... -v passes all 26 tests, go vet and gofmt clean. With a fake runner, TestAnswerQuestionSendsSessionMessage types 'Use sqlite' through the real rexuikit Form TextArea and sees session.message [S1, Use sqlite], then tea.Quit; TestAnswerBlockerCallsQuestionAnswer sees question.answer [which-db, C5, sqlite], then Quit; TestAnswerEscAndEmptyCallNothing sees an empty Enter and a cancel make no call past the inbox load. The attempt-2 fix (send() reopens a fresh answer form via answerForm on an empty answer or a failed send, since a submitted Form drops keys) is covered: TestAnswerEmptyThenAnswerSends shows the refusal 'an answer needs words' with no call, then typing into the new form sends and quits; TestAnswerSendFailureReopensForm shows 'send failed: no session', and a retry in the reopened form sends session.message [S1, Use pg] and quits. Read the diff: it touches only send/openAnswer in internal/tui/inboxapp.

## Links

### Out

- verifies [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
