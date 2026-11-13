# queue-wait-journal-reader

capability · `01M79DE919PHFX4JVJ8V97PGE5` · parent [daemon-picks](../daemon-picks.md)

## acceptance

go test -trimpath ./internal/daemon/... -run QueueWait: appended lines read back oldest first with every field; a malformed line between two good ones is skipped; a missing file reads as no lines; lines appended from several goroutines at once each read back whole

## statement

internal/daemon defines the queue-wait journal's line (at, change, title, reason, detail), appends one to queue.wait.jsonl in the daemon's folder under the git common dir with a single write per line, and reads the file back as every line oldest first, skipping lines that do not parse

## Links

### Out

- decided-by [queue-wait-journal-file](../../../../../decision/queue-wait-journal-file.md)

### In

- depends-on [queue-wait-journal-lines](queue-wait-journal-lines.md)
- verifies [queue-wait-journal-reader-verified](../../../../../evidence/queue-wait-journal-reader-verified.md)
- verifies [queue-wait-journal-reader-qa-1](../../../../../evidence/queue-wait-journal-reader-qa-1.md)
- verifies [queue-wait-journal-reader-qa-2](../../../../../evidence/queue-wait-journal-reader-qa-2.md)
- verifies [pr-536-queue-wait-journal](../../../../../evidence/pr-536-queue-wait-journal.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
