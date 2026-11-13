# queue-wait-journal-lines-verified

evidence · `01M79VZA5J415GGMEVHF50RDF8` · actor [claude-code-18654100](../actor/claude-code-18654100.md)

## result

PASS

## statement

At c0b442e6a re-ran the Checks block: go test -trimpath -count=1 ./internal/api/... -run QueueWait (TestQueueWaitJournal, TestQueueWaitFirstPollStarts, TestQueueWaitRestartWritesNoDuplicate all PASS, exit 0), go test -trimpath ./internal/daemon/... (ok, exit 0), rhiza ship generate (exit 0), rhiza render (exit 0, tree unchanged). Read the diff: daemonPoll1 journals the poll's titled skips after every skip source is gathered and a started line after each daemonStart; JournalQueueWaits writes only pre-start reasons (the 11 listed), skips change sets that RunStates.PostStart reports live/held/shipped/failed (shared with Pick's decide), and dedups against the last reason read from the journal file, so a restarted daemon writes no duplicate. Tests cover no-slot x2, usage-reserve, load-high under a machine hold, start, then held and live polls under every hold writing nothing (exactly four lines with change, title and detail), first-poll start writing only started, and a fresh daemon writing no duplicate. Diff stays within internal/api, internal/daemon and the builder's actor record; no stubs or TODOs.

## Links

### Out

- verifies [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
