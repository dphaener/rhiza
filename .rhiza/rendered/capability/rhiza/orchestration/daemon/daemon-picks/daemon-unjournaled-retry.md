# daemon-unjournaled-retry

capability · `01M6A0NNP5H447ACNTCPMS7D9F` · parent [daemon-picks](../daemon-picks.md)

## acceptance

TestDaemonHeld in internal/api/daemon_picks_test.go: an unjournaled change set is still held at 15 minutes less a second, is released at 15 minutes, and is pending again after it is picked; a landed change set leaves the hold with no new journal; daemon-status.txtar's Not journaled line names when it will be retried

## statement

A change set whose daemon start wrote no journal is held for 15 minutes and then picked again, without a daemon restart; one that has landed or been withdrawn is dropped from the hold at once, and the one-minute grace before a start counts as unjournaled is unchanged

## Links

### Out

- decided-by [unjournaled-block-expires](../../../../../decision/unjournaled-block-expires.md)

### In

- verifies [daemon-unjournaled-retry-verified](../../../../../evidence/daemon-unjournaled-retry-verified.md)
- verifies [daemon-unjournaled-retry-qa-1](../../../../../evidence/daemon-unjournaled-retry-qa-1.md)
- verifies [pr-154-unjournaled-hold-expires](../../../../../evidence/pr-154-unjournaled-hold-expires.md)

## History

- `01M6A0JVE5ZV75XEMB0ZNYN8SR` A daemon start that wrote no journal is retried after a hold, without restarting the daemon
