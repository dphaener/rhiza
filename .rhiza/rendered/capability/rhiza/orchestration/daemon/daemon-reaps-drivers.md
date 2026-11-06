# daemon-reaps-drivers

capability · `01M6Q43V93A114TYBBXS89ZV6Z` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestReap' on Linux: a test starts a Setsid child recorded on the driver list that nothing waits on, lets it exit, runs the reaper, and finds the pid gone with no zombie left and off the list, while a plain os/exec child that calls setsid() itself, running at the same time, still gets its own exit status from Wait; the reaper runs at daemon start and on each tick, and does nothing off Linux

## statement

The daemon reaps its finished run drivers, including those started by an image before a re-exec, so none stays a zombie: each detached driver's pid goes on a driver list under the git common dir, and at start and on a timer the daemon waits by pid, without blocking, only on listed pids that are still its children, so the git, make or hook children os/exec waits on are never reaped by it, even one that calls setsid itself

## Links

### Out

- decided-by [reap-own-session-children](../../../../decision/reap-own-session-children.md)

### In

- verifies [daemon-reaps-drivers-verified](../../../../evidence/daemon-reaps-drivers-verified.md)
- verifies [daemon-reaps-drivers-qa-1](../../../../evidence/daemon-reaps-drivers-qa-1.md)
- verifies [daemon-reaps-drivers-verified-2](../../../../evidence/daemon-reaps-drivers-verified-2.md)
- verifies [daemon-reaps-drivers-qa-2](../../../../evidence/daemon-reaps-drivers-qa-2.md)
- verifies [pr-321-fix-the-daemon-never-reaps-its](../../../../evidence/pr-321-fix-the-daemon-never-reaps-its.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
