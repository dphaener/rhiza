# inbox-app-answer-verified

evidence · `01M6NAKDSXFD7XQQGJE4K9QGCG` · actor [claude-code-7b679c48](../actor/claude-code-7b679c48.md)

## result

PASS

## statement

At 7bfe3f64 ran go test -trimpath -count=1 -v ./internal/tui/inboxapp/...: all 12 tests pass, plus go vet and gofmt clean. With the fake runner, TestAnswerQuestionSendsSessionMessage shows Enter on a question opens a Form with one TextArea ("Your answer"); submitting "Use sqlite" calls session.message [S1, Use sqlite] and returns tea.Quit. TestAnswerBlockerCallsQuestionAnswer shows a blocker submits question.answer [which-db, C5, sqlite] and quits. TestAnswerEscAndEmptyCallNothing shows an empty submit neither quits nor calls anything, and cancelling quits with only the initial inbox load call. Read the diff: Update quits on FormCancelledMsg and on an esc key while the form shows, before any run call; send trims the text and returns nil when it is empty; cmd/rhiza/tui.go maps the three-arg question.answer call to --change plus question and text, matching QuestionAnswerRequest.

## Links

### Out

- verifies [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
