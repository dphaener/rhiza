# queue-wait-journal-lines-verified-2

evidence · `01M79X840ZDSKXR877F4HBFV16` · actor [claude-code-cf22de14](../actor/claude-code-cf22de14.md)

## result

PASS

## statement

Verified attempt 2 at 6f0d5d935. Re-ran every Checks: command, each exit 0: go test -trimpath ./internal/api/... -run QueueWait (TestQueueWaitJournal, TestQueueWaitFirstPollStarts, TestQueueWaitRestartWritesNoDuplicate, TestQueueWaitShipOwnerStartedOutside all PASS); go test -trimpath ./internal/daemon/...; go vet ./internal/daemon/ ./internal/api/; rhiza ship generate; rhiza render (no tracked diff after); make self-check (only pre-existing unclaimed warnings). Read the diff: JournalQueueWaits now takes the run journals and, before the skip filters, appends a started line at the run's start time with detail 'started outside the daemon' for a change set whose last line is a pre-start reason and whose run started after it, as Decision queue-wait-journal-outside-starts says. The new test covers ship-owner, then a run started by hand, then started at the run's time with that detail, and no further line on a later poll. The earlier acceptance points still pass. The leaf's statement and acceptance edits match the decision. No TODOs or stubs, and nothing changed outside the leaf beyond the builder's actor record and the rendered files.

## Links

### Out

- verifies [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
