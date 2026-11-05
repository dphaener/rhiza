# run-send-back-qa-2

evidence · `01M6NDH68VZE1RXZC37XCSKEED` · actor [claude-code-3293f9c2](../actor/claude-code-3293f9c2.md)

## result

PASS

## statement

QA at 433db495: send-back from the built inbox app on a real review run waiting at signoff journalled the note as a finding; rhiza run status lists it under Sign-offs as 'changes asked by test' with 'alpha f1 major at owner: tighten the tests', and no approval is journalled. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-send-back|inbox-signoff)$' passes.

## Links

### Out

- verifies [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
