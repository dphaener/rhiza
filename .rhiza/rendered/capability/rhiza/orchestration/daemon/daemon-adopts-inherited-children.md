# daemon-adopts-inherited-children

capability · `01M7A01JB2302EDM0R0N111NEV` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestReap' on Linux: a test starts a plain child it never waits on (standing for one inherited through a re-exec), runs the start-time adoption, lets the child exit, runs the reaper and finds the pid gone with no zombie and off the list, while a child os/exec starts after the adoption still gets its own exit status from Wait; off Linux adoption does nothing. Touches cmd/rhiza/drivers_linux.go, cmd/rhiza/drivers_unix.go, cmd/rhiza/drivers_other.go, cmd/rhiza/detach_other.go, cmd/rhiza/reap_test.go, cmd/rhiza/ctx.go, internal/api/api.go, internal/api/daemon.go

## statement

When a daemon image starts, on Linux, it puts every child process it already has (one an earlier image left running when it re-executed, whose waiter died with that image) on the home project's driver list as its own, before it starts any program, so the reaper waits on each by pid and none stays a defunct [git]; a child the new image starts later is never put on the list by this

## Links

### Out

- decided-by [decision-adopt-inherited-children](../../../../decision/decision-adopt-inherited-children.md)

### In

- verifies [daemon-adopts-inherited-children-verified](../../../../evidence/daemon-adopts-inherited-children-verified.md)
- verifies [daemon-adopts-inherited-children-qa-1](../../../../evidence/daemon-adopts-inherited-children-qa-1.md)
- verifies [daemon-adopts-inherited-children-qa-2](../../../../evidence/daemon-adopts-inherited-children-qa-2.md)
- verifies [daemon-adopts-inherited-children-qa-3](../../../../evidence/daemon-adopts-inherited-children-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
