# inbox-app-retry-abandon-qa-1

evidence · `01M6NCQ5J5JY75FXW56EKSA7NQ` · actor [claude-code-f41a81b4](../actor/claude-code-f41a81b4.md)

## result

PASS

## statement

QA at e4777542 with the built binary in tmux: Enter on a failed run opens a form with a Retry/Abandon Select. Abandon then No re-showed the choice ('abandon needs the confirm') with the journal unchanged; Esc quit with nothing done; Abandon then Yes appended finished 'abandoned by t' and exited 0. On a real failed shiptest run (check session exiting 1), Retry queued it in continues.json with retry:true and the app exited 0.

## Links

### Out

- verifies [inbox-app-retry-abandon](../capability/rhiza/inbox-app-retry-abandon.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
