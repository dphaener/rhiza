# question-ab537112-1

question · `01M6Q41DABCQW2A1P9TRH1Y2E6`

## session

ab537112-dbb7-4d84-807c-01a1ae3eb5c1

## statement

How should the daemon reap its finished run drivers? (Fix The daemon never reaps its, 01M6Q3X1.) The intent proposes a periodic wait4(-1, WNOHANG), but the daemon also runs git, make and hook commands through os/exec, some concurrently from web UI request goroutines; wait4(-1) can reap one of those first, so its cmd.Wait/Output fails with 'no child processes'. Detached drivers are the only children started with Setsid, so they lead their own session.

## Links

### In

- answers [reap-own-session-children](../decision/reap-own-session-children.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
