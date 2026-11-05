# m26-inbox-app

intent · `01M6MZ0BVXSHHDCB4GXVABYTKE` · parent [milestone-26](../milestone-26.md)

## statement

From the Inbox app the owner picks an inbox item and does its one action (answers a question or gate, signs off a run, retries or abandons a failed run), and the app closes.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-run-actions](../../decision/tui-run-actions.md)

### In

- serves change set `01M6MYN9GRWHRJSHA31D9S838H` (Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces)

## History

- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
