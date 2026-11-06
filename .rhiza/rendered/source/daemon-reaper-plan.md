# daemon-reaper-plan

source · `01M6Q43TY2EPDME3KRYNNZ66MQ`

## quotes

```json
[]
```

## raw_sha256

52e4a8637950c6bb72c01ae92930108b811801d1a4739204473603168b973a64

## session

ab537112-dbb7-4d84-807c-01a1ae3eb5c1

## summary

Autonomous plan of the fix for the daemon leaving zombie run drivers. The code already reaps each detached driver with go cmd.Wait(), so the zombies come from in-place re-execs, after which the new image does not know its children. The planner raised that the intent's wait4(-1, WNOHANG) would also reap the git, make and hook children os/exec is waiting on (some from concurrent web UI requests). The arbiter first decided the daemon reaps only its own children that lead their own session (detached drivers), at start and on a timer, by pid, on Linux via /proc. Plan review round 1 (f1) showed session leadership is not enough, since a hook running setsid(1) makes the daemon's direct child a session leader, so the decision was revised: each detached driver's pid is recorded on a driver list under the git common dir, which outlives a re-exec, and the reaper, at start and on a timer, waits (wait4(pid, WNOHANG)) only on listed pids that are still the daemon's children, on Linux via /proc. One leaf under daemon, served by milestone 25 like the other fixes.

## title

Daemon reaping its run drivers, autonomous plan

## Links

### In

- derived-from [reap-own-session-children](../decision/reap-own-session-children.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
