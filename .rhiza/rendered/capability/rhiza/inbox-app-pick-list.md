# inbox-app-pick-list

capability · `01M6MZ9RZ8M65B2FB58A75WYZ2` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: a row per inbox item in order; Enter on a question, gate, signoff, failed or accept item opens that item's action; Enter on a stalled, overlap, disk-low, usage-reserve or main-stale item shows its detail and asks the runner to open its run's tab; an empty inbox shows one empty-state line; the tui-inbox-app testscript compares rhiza tui app inbox --once with rhiza inbox --json byte for byte.

## statement

The Inbox app lists every inbox item as one rexuikit List, a row per item with its kind, change set title and what it waits for; Enter picks the item and opens its action (answer for questions and gates, sign-off, retry or abandon for a failed run, accept for an accept item), and any other kind shows its detail and opens the run's tab; --once prints what rhiza inbox --json prints.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)
- decided-by [tui-inbox-all-items](../../decision/tui-inbox-all-items.md)

### In

- blocks [inbox-app-kinds](../../question/inbox-app-kinds.md)
- verifies [inbox-app-pick-list-verified](../../evidence/inbox-app-pick-list-verified.md)
- verifies [inbox-app-pick-list-qa-1](../../evidence/inbox-app-pick-list-qa-1.md)
- verifies [inbox-app-pick-list-verified-2](../../evidence/inbox-app-pick-list-verified-2.md)
- verifies [inbox-app-pick-list-qa-2](../../evidence/inbox-app-pick-list-qa-2.md)
- verifies [inbox-app-pick-list-verified-3](../../evidence/inbox-app-pick-list-verified-3.md)
- verifies [inbox-app-pick-list-qa-3](../../evidence/inbox-app-pick-list-qa-3.md)
- verifies [pr-299-inbox-app-on-the-modal-shell](../../evidence/pr-299-inbox-app-on-the-modal-shell.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
