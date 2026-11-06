# daemon-reaps-drivers-verified

evidence · `01M6Q7KXVZKE0CQKD7YCRJECZT` · actor [claude-code-e0dfc339](../actor/claude-code-e0dfc339.md)

## result

PASS

## statement

At fd4fdfe0 on Linux, ran go test -trimpath -count=1 -v ./cmd/rhiza -run TestReap: PASS. It was stable with -count=3. The test starts a Setsid child put on the driver list with addDriver that nothing waits on, waits until /proc shows it as a zombie, runs reapDrivers(common), and finds /proc/<pid> gone and the list empty. A concurrent unlisted os/exec child with Setsid still gets exit status 3 from Wait. Read the diff. detachProgram adds each driver pid to <common>/rhiza/drivers, which is beside rhiza/logs and the same path that reapDrivers reads. reapList calls Wait4(pid, WNOHANG) only on listed pids whose /proc ppid is the daemon, and drops pids that are reaped or not its children, so unlisted git, make or hook children are never touched. daemonLoop calls ctx.Reap(common) at the top of each iteration: once at start and once after each tick. newCtx wires Reap to reapDrivers. Off Linux, reapList and reapDrivers are no-ops. go vet passed, GOOS=darwin and GOOS=windows builds passed, and go test ./internal/api -run Daemon passed.

## Links

### Out

- verifies [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
