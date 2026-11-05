# tui-inbox-all-items

decision · `01M6N0763RMWPEPXWFA3RQ8CB1`

## statement

The Inbox app lists every inbox item: questions and gates are answered, a sign-off is signed off, a failed run is retried or abandoned, an accept item is accepted; every other kind (stalled, overlap, disk-low, usage-reserve, main-stale) shows its detail and opens the run's tab, with no special action

## wording

All items, act or open: questions and gates answer, sign-off signs off, failed runs retry or abandon, accept items accept, and the rest show their detail and open the run's tab, the agent's option, picked by the owner

## Links

### Out

- answers [inbox-app-kinds](../question/inbox-app-kinds.md)
- derived-from [tui-small-apps-interview](../source/tui-small-apps-interview.md)

### In

- decided-by [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
