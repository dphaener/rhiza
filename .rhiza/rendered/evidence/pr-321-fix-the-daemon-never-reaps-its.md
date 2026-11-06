# pr-321-fix-the-daemon-never-reaps-its

evidence · `01M6QCPA5TW2X385TXMYHKNYX1` · actor [dphaener](../actor/dphaener.md)

## head

819c74a4a33e7fbb9f643c8524ede9194740a27a

## pr

https://github.com/dphaener/rhiza/pull/321

## result

1 check passed: check SUCCESS

## statement

Pull request #321 (https://github.com/dphaener/rhiza/pull/321) passed CI at 819c74a

## status

verified

## Links

### Out

- verifies [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
