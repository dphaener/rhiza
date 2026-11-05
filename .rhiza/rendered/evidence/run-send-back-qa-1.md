# run-send-back-qa-1

evidence · `01M6NCQC2W6YE7PY6G0YX39DHW` · actor [claude-code-f41a81b4](../actor/claude-code-f41a81b4.md)

## result

PASS

## statement

QA at e4777542 with the built binary on a scratch repo whose shiptest run (fake claude/gh/rex) stopped at signoff: send-back through the inbox app with note 'say more' journalled signoff outcome changes with one finding (alpha, major, place owner, reason 'say more'), rhiza run status lists it under Sign-offs as changes asked, the run is queued for the daemon and no approved signoff is journalled. go test -trimpath ./cmd/rhiza -run TestScripts/run-send-back passes.

## Links

### Out

- verifies [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
