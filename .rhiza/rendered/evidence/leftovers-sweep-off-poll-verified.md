# leftovers-sweep-off-poll-verified

evidence · `01M7A2BKZG3QEWM7QMNZ31FX7T` · actor [claude-code-aabdfee0](../actor/claude-code-aabdfee0.md)

## result

PASS

## statement

At 0f54a0aea re-ran the Checks block: go test -trimpath ./internal/api -run TestDaemonSweepLeftovers (exit 0; TestDaemonSweepLeftoversOffPoll and TestDaemonSweepLeftoversPrunesOnlyOrphans pass), go test -trimpath ./internal/daemon -run Sweep (exit 0, incl. TestSweepScheduleStartOnce), rhiza ship generate (exit 0, nothing changed), rhiza render (exit 0, tree clean). Also -race -count=3 on both packages passed, gofmt and go vet clean. Diff read: daemonSweepLeftovers (called from the home project's atHome poll step in daemon.go) claims the sweep with SweepSchedule.Start and runs it on a goroutine over a copy of ctx and commons; Start refuses while one runs; the stderr line text is unchanged. The new test blocks the sweep on a channel, asserts the step returns, that a due poll two hours on starts no second sweep, and that the freed line is written once after release. No stubs or TODOs; changes stay in the leaf's listed files plus the builder's actor record and rendered files.

## Links

### Out

- verifies [leftovers-sweep-off-poll](../capability/rhiza/orchestration/daemon/leftovers-sweep-off-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
