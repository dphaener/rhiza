# daemon-reaps-drivers-qa-2

evidence · `01M6QAMEHF00N059R0B5NJMYTF` · actor [claude-code-c6ad283b](../actor/claude-code-c6ad283b.md)

## result

PASS

## statement

At 32103a7d on Linux, built the binary (go build -trimpath) and ran it as a real daemon in a scratch git repo after rhiza init, with a scratch HOME. sh started four 'sleep 3' children and wrote .git/rhiza/drivers with three lines: A recorded under its own pid and its /proc start time (field 22), B recorded under its own pid with the wrong start time '1', and D as a bare legacy pid. C was left off the list. sh then exec'd into 'rhiza daemon --min-free 0 --permission-mode acceptEdits', keeping the pid and start time, as a re-exec in place does. The reap at daemon start dropped B and D from the list without waiting on them and kept A, which was still running. By t=5 s A had been reaped (gone from ps), and the list file was removed. B, C and D stayed zombies whose ppid was the daemon, so entries from another daemon identity are never waited on and unlisted children are untouched. They disappeared once the daemon was stopped with SIGINT. go test -trimpath -count=1 -v ./cmd/rhiza -run TestReap passes TestReap and TestReapStale. GOOS=darwin go vet ./cmd/rhiza is clean. Not exercised: a driver started through the daemon's own detachProgram (a real run start), which TestReap covers through addDriver.

## Links

### Out

- verifies [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
