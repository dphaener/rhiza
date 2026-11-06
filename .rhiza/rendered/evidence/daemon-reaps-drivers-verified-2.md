# daemon-reaps-drivers-verified-2

evidence · `01M6QAFPEXBCDAPPM9V9CQKQJ1` · actor [claude-code-ba3d8b7b](../actor/claude-code-ba3d8b7b.md)

## result

PASS

## statement

At 92ec88e8 on Linux, go test -trimpath -count=1 -v ./cmd/rhiza -run TestReap passes TestReap and TestReapStale. TestReap starts a Setsid true child, records it with addDriver, waits until it is a zombie, runs reapDrivers, then finds /proc/<pid> gone and the list empty, while a concurrent plain os/exec child running sh with Setsid still gets exit status 3 from Wait. TestReapStale records a live, unwaited child under a wrong daemon start time and under a wrong daemon pid. After reapDrivers the child is still a zombie (not reaped), the list is empty (dropped), and its Wait returns exit 5. This covers decision reap-only-own-image-entries. In the diff, entries are '<driver pid> <daemon pid> <daemon start>' with start from field 22 of /proc/<pid>/stat. reapList calls wait4(pid, WNOHANG) only when the entry's daemonID matches selfID, start is non-empty and ppid is self, under driverMu, and treats ECHILD as dropped. internal/api daemonLoop calls ctx.Reap before each poll, so at start and on each tick. Off Linux, drivers_other.go (unix, not linux) and detach_other.go have no-op reapers, and GOOS=darwin go vet ./cmd/rhiza is clean. go vet ./cmd/rhiza is clean, and internal/daemon tests pass. internal/api passes with TestReadBudget skipped. TestReadBudget, a load-sensitive timing test that does not touch Reap, hung under load average ~21.

## Links

### Out

- verifies [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
