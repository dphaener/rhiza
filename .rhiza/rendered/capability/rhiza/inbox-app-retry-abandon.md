# inbox-app-retry-abandon

capability · `01M6MZAFPVHK2JM4JJAP4953HC` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: retry calls run.resume and abandon calls run.abandon on the item's run, each then quitting the app; abandon without the Confirm and Esc call nothing; neither choice is offered on other kinds.

## statement

Picking a failed, gave-up or skipped run offers retry or abandon in a rexuikit Form (a Select, with a Confirm before abandoning); retry queues run.resume, abandon runs run.abandon, and the app closes after either.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [inbox-app-retry-abandon-verified](../../evidence/inbox-app-retry-abandon-verified.md)
- verifies [inbox-app-retry-abandon-qa-1](../../evidence/inbox-app-retry-abandon-qa-1.md)
- verifies [inbox-app-retry-abandon-qa-2](../../evidence/inbox-app-retry-abandon-qa-2.md)
- verifies [pr-299-inbox-app-on-the-modal-shell](../../evidence/pr-299-inbox-app-on-the-modal-shell.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
