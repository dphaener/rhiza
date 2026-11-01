# inbox-disk-low-qa-1

evidence · `01M6B38J1ZFQ6MYY39NAMKVPE3` · actor [claude-code-e1013c3f](../actor/claude-code-e1013c3f.md)

## result

PASS

## statement

At 584cd6c3, in the same scratch repo with the built binary: after a real low-disk poll (--min-free 100G), rhiza inbox lists one '<id> disk-low' item reading '57.8 GiB free, under the minimum of 100.0 GiB: the daemon holds new starts and resumes'. The ID stays the same across a second low poll. After a poll with --min-free 1G it prints 'Nothing waits on you'. The inbox app view was checked only through go test ./internal/tui/inboxapp (PASS), not by hand

## Links

### Out

- verifies [inbox-disk-low](../capability/rhiza/orchestration/inbox/inbox-list/inbox-disk-low.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
