# tui-signoff-approve-or-send-back

decision · `01M6N07EQSBWX87XGQZ08K60ER`

## statement

The Inbox app's sign-off dialog offers Approve, and Send back with a note; the note becomes findings for another fixer cycle of the run

## wording

Approve or send back: the dialog has Approve, and Send back with a note that becomes findings for another fixer cycle, the agent's option, picked by the owner

## Links

### Out

- answers [inbox-signoff-request-changes](../question/inbox-signoff-request-changes.md)
- derived-from [tui-small-apps-interview](../source/tui-small-apps-interview.md)

### In

- decided-by [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)
- decided-by [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
