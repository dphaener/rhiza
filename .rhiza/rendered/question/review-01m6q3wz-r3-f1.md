# review-01m6q3wz-r3-f1

question · `01M6Q4N00EHQZ53Z49E4AGMB9J`

## statement

Plan review round 3 left finding f1 (contradiction on reap-own-session-children): The decision says the only limit left is pid reuse in the instant between a goroutine's Wait and dropping the pid, and that os/exec children (git, make, hooks) are never on the list. But the list lives under the git common dir and outlives not only a re-exec but a full daemon restart (stop/start, crash, systemd restart) with a new pid. Drivers of the old daemon are then reparented and reaped by init, yet their pids stay on the list; the new daemon only drops an entry once it scans it and sees ppid != its own pid. If any listed pid has meanwhile been reused by a git/make/hook child of the new daemon (ppid == new daemon), the reaper's wait4(pid, WNOHANG) can reap it and that cmd.Wait fails with 'no child processes': the very race the decision rules out. 'Runs once when the daemon starts' does not say it runs before the daemon starts any child, so the startup scan does not close it. Fix: record the daemon's pid with each entry (a re-exec keeps the pid) and reap only entries recorded by the running pid, dropping the rest; or clear the list on a fresh (non-re-exec) start before any child is started; and cover the stale-entry case in the leaf's test. The critic confirmed it: reap-own-session-children keeps the driver list under the git common dir and drops an entry only when a scan finds its ppid is not the daemon's pid. A fresh daemon start (crash, stop/start; internal/daemon/lock.go lets a new holder take over with a new pid, unlike reexec_unix.go's same-pid syscall.Exec) leaves the old daemon's driver pids on the list after init reaps them. Neither the decision nor the daemon-reaps-drivers leaf says the startup scan runs before the daemon starts any child, and the timer scan does not care which image wrote an entry. So a git/make/hook child of the new daemon that draws a stale listed pid gets ppid == daemon and is reaped by wait4(pid, WNOHANG), and its cmd.Wait fails with ECHILD. That contradicts the decision's claim that the only limit left is the Wait-to-drop instant, and the leaf's acceptance check does not cover stale entries. Fix: tag each entry with the recording daemon's pid (a re-exec keeps it) and reap only entries matching the running pid, or clear the list on a non-re-exec start before any child starts; add the stale-entry case to the test. How should it be settled?

## Links

### In

- answers [reap-only-own-image-entries](../decision/reap-only-own-image-entries.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
