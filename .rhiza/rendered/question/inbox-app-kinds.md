# inbox-app-kinds

question · `01M6MZANYQFS6TDEKFS3P4RT6D`

## statement

Which inbox kinds does the Inbox app list and act on? The owner named answering a question or gate, signing off and retrying or abandoning a failed run: is a gate a blocker (a question blocking a change set), an accept item, or both, and are the other kinds (open, stalled, overlap, disk-low, usage-reserve, main-stale) listed read-only, left out, or given an action? Do the old app's other actions stay (open the item's Rex tab, take the advisor's recommendation, nudge a stalled session)?

## Links

### Out

- blocks [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

### In

- answers [tui-inbox-all-items](../decision/tui-inbox-all-items.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
