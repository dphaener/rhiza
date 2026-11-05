# inbox-signoff-request-changes

question · `01M6MZATM6JDM5M8P24TCH24JM`

## statement

Does signing off a run in the Inbox app only approve it, or can the owner also request changes there (inbox.request-changes, as the old app's c did)?

## Links

### Out

- blocks [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)

### In

- answers [tui-signoff-approve-or-send-back](../decision/tui-signoff-approve-or-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
