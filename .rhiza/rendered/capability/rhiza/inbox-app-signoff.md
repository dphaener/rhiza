# inbox-app-signoff

capability · `01M6MZA8WT3TB8HWR4D33QMMQ5` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: Approve on a signoff item calls run.approve with its run and the app quits; Send back with a note calls run.send-back with the run and the note and quits; Send back with an empty note is refused; Esc calls nothing.

## statement

Picking a sign-off opens a rexuikit Dialog naming the run and what it waits for, with Approve and Send back: Approve signs it off through run.approve and Send back takes a note through run.send-back, either then closing the app; Esc closes it with nothing done.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)
- decided-by [tui-signoff-approve-or-send-back](../../decision/tui-signoff-approve-or-send-back.md)

### In

- blocks [inbox-signoff-request-changes](../../question/inbox-signoff-request-changes.md)
- verifies [inbox-app-signoff-verified](../../evidence/inbox-app-signoff-verified.md)
- verifies [inbox-app-signoff-qa-1](../../evidence/inbox-app-signoff-qa-1.md)
- verifies [inbox-app-signoff-qa-2](../../evidence/inbox-app-signoff-qa-2.md)
- verifies [pr-299-inbox-app-on-the-modal-shell](../../evidence/pr-299-inbox-app-on-the-modal-shell.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
