# queue-wait-journal-lines

capability · `01M79DE92G4FFVHC5MFPY6EFGN` · parent [daemon-picks](../daemon-picks.md)

## acceptance

go test -trimpath ./internal/api/... -run QueueWait: polls that skip a change set for no-slot twice, then for usage-reserve, then for load-high under a machine hold, then start it, then skip it as held and then live journal exactly four lines (no-slot, usage-reserve, load-high, started) with its change, title and detail; a change set started on its first poll journals only a started line; a fresh daemon over the same journal writes no line for an unchanged reason; a change set skipped as ship-owner whose run is then started outside the daemon journals ship-owner then started, the started line's at being the run journal's start time and its detail 'started outside the daemon'

## statement

Each daemon poll appends a queue-wait line for an accepted change set only when its pre-start wait changes: when a poll first skips it for a pre-start reason, with the poll's Skip reason and detail (no-slot, not-ready, no-milestone, no-leaves, disk-low, load-high, memory-low, usage-reserve, paused, ship-owner and workflow-unresolved), when that reason changes, and a started line when its run starts. A run the daemon starts is journaled as started when the daemon starts it. A run the daemon did not start (one started by hand, or a ship-owner change set's) is journaled on the next poll, before any pre-start filter: when a run journal under <common>/rhiza/runs/ has a started entry for a change set later than its last line and that line is a pre-start reason, the poll appends a started line whose at is the run's start time and whose detail is 'started outside the daemon'. The post-start Skip reasons (live, held, unjournaled, failed and shipped, which the daemon gives a change set whose run has already started: internal/daemon/pick.go, internal/api/daemon.go) are never journaled, so after started nothing more is written for that change set until a poll skips it for a pre-start reason again (a retried or reopened change set waiting anew). The last reason per change set is taken from the journal, so a restarted daemon writes no duplicate line

## Links

### Out

- depends-on [queue-wait-journal-reader](queue-wait-journal-reader.md)
- decided-by [queue-wait-journal-file](../../../../../decision/queue-wait-journal-file.md)

### In

- blocks [review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f1](../../../../../question/review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f1.md)
- verifies [queue-wait-journal-lines-verified](../../../../../evidence/queue-wait-journal-lines-verified.md)
- verifies [queue-wait-journal-lines-qa-1](../../../../../evidence/queue-wait-journal-lines-qa-1.md)
- verifies [queue-wait-journal-lines-verified-2](../../../../../evidence/queue-wait-journal-lines-verified-2.md)
- verifies [queue-wait-journal-lines-qa-2](../../../../../evidence/queue-wait-journal-lines-qa-2.md)
- verifies [pr-536-queue-wait-journal](../../../../../evidence/pr-536-queue-wait-journal.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
