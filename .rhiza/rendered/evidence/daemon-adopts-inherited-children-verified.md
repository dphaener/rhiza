# daemon-adopts-inherited-children-verified

evidence · `01M7A502V4H7B8HM4QDV07WQ8Q` · actor [claude-code-38a90e3a](../actor/claude-code-38a90e3a.md)

## result

PASS

## statement

At d392015b8 on Linux: go test -trimpath ./cmd/rhiza -run TestReap -v passed TestReap, TestReapStale and TestReapAdopt (exit 0); go test -trimpath ./internal/api -run Daemon exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only this verifier's actor record written). TestReapAdopt starts an unwaited child, runs adoptChildren, finds it listed as selfID, lets it exit, reapDrivers removes the pid with no zombie and empties the list, while a child started after adoption is not listed and its Wait returns exit status 4. Off Linux adoptList/adoptChildren are no-ops (drivers_other.go, detach_other.go; TestReapAdoptOther in reap_other_test.go); GOOS=darwin test build and GOOS=windows vet of ./cmd/rhiza pass. Diff read: runDaemon calls ctx.Adopt(common) once, after taking the lock and before any program starts; Ctx.Adopt wired in ctx.go; no stubs or TODOs; only the listed files plus the off-Linux test file and the builder's actor record changed.

## Links

### Out

- verifies [daemon-adopts-inherited-children](../capability/rhiza/orchestration/daemon/daemon-adopts-inherited-children.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
