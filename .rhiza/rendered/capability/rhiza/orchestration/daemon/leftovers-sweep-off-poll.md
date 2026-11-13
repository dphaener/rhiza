# leftovers-sweep-off-poll

capability · `01M7A01HSS3X0K4M3EM9XWRDFJ` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/api -run 'TestDaemonSweepLeftovers': a test makes the sweep block on a channel, runs the home project's poll step and asserts it returns while the sweep is still blocked, that a second due poll starts no second sweep, and that the freed line is written once the sweep is released; the existing orphan-pruning test still passes. Touches internal/api/daemon.go, internal/api/daemon_leftovers.go, internal/api/daemon_leftovers_test.go, internal/daemon/sweepsched.go, internal/daemon/sweepsched_test.go

## statement

The daemon runs its hourly leftovers sweep on a goroutine of its own instead of inside the poll: a poll that finds the sweep due starts it and goes on, a poll never starts a second sweep while one runs, and the sweep's log line is unchanged

## Links

### Out

- depends-on [scratch-sweep-one-listing](scratch-sweep-one-listing.md)
- decided-by [decision-housekeeping-off-poll](../../../../decision/decision-housekeeping-off-poll.md)

### In

- verifies [leftovers-sweep-off-poll-verified](../../../../evidence/leftovers-sweep-off-poll-verified.md)
- verifies [leftovers-sweep-off-poll-qa-1](../../../../evidence/leftovers-sweep-off-poll-qa-1.md)
- verifies [leftovers-sweep-off-poll-verified-2](../../../../evidence/leftovers-sweep-off-poll-verified-2.md)
- verifies [leftovers-sweep-off-poll-qa-2](../../../../evidence/leftovers-sweep-off-poll-qa-2.md)
- verifies [leftovers-sweep-off-poll-qa-3](../../../../evidence/leftovers-sweep-off-poll-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
