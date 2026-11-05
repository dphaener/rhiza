# inbox-app-answer

capability · `01M6MZA0N04GRKJ6486V5XRAX9` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: submitting on a question calls session.message with the item's session and the typed text, on a blocker calls question.answer with the question, its change set and the typed text, then the app quits; Esc or an empty answer calls nothing.

## statement

Picking a question or a blocker opens a rexuikit Form with one TextArea; submitting sends the owner's typed words (session.message to the asking session for a run's question, question.answer for a blocker) and closes the app, and Esc closes it with nothing sent.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [inbox-app-answer-verified](../../evidence/inbox-app-answer-verified.md)
- verifies [inbox-app-answer-qa-1](../../evidence/inbox-app-answer-qa-1.md)
- verifies [inbox-app-answer-verified-2](../../evidence/inbox-app-answer-verified-2.md)
- verifies [inbox-app-answer-qa-2](../../evidence/inbox-app-answer-qa-2.md)
- verifies [pr-299-inbox-app-on-the-modal-shell](../../evidence/pr-299-inbox-app-on-the-modal-shell.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
