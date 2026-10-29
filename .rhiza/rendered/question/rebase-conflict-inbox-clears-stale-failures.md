# rebase-conflict-inbox-clears-stale-failures

question · `01M6350ZQZZQP18DGDGDWQGP0G`

## statement

Rebasing branch cs/inbox-clears-stale-failures onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/daemon.go, internal/api/inbox.go?

## Links

### In

- answers [inbox-stale-rebase-resolution](../decision/inbox-stale-rebase-resolution.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
