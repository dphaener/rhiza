# daemon-reaps-drivers-qa-1

evidence · `01M6Q8ERFZVZH712M7WDHYM5K0` · actor [claude-code-c745c78d](../actor/claude-code-c745c78d.md)

## result

PASS

## statement

At 9a3773fa on Linux, built the binary (go build -trimpath) and ran it as a real daemon in a scratch git repo with rhiza init: sh started two 'sleep 3' children, wrote the first one's pid to .git/rhiza/drivers, then exec'd into 'rhiza daemon --min-free 0 --permission-mode acceptEdits', the same pid, as a re-exec does. After 3 s, ps showed both sleeps as zombies (Z) whose ppid was the daemon. The listed one disappeared on the daemon's next loop pass, at about 72 s, because the first poll is slow here (its usage scan of ~/.claude transcripts predates this change). The list file was then removed. The unlisted sleep stayed a zombie of the daemon until the daemon stopped, so the reaper waits only on listed pids and never calls wait4(-1). go test -trimpath -count=1 ./cmd/rhiza -run TestReap: ok. Not exercised: stale entries left by a different daemon, which the code cannot tell apart (see the qa finding).

## Links

### Out

- verifies [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
