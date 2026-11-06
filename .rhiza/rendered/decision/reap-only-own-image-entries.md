# reap-only-own-image-entries

decision · `01M6Q4P92SS3KB3ASVQ8X7HZPT`

## reasoning

The finding is right: reap-own-session-children keeps the list under the git common dir and drops an entry only when its ppid is not the daemon's, so after a fresh daemon start (internal/daemon/lock.go lets a new holder take over with a new pid, unlike reexec_unix.go's same-pid syscall.Exec) stale pids stay listed, and one reused by a new os/exec child would be reaped and its cmd.Wait would fail with ECHILD. Of the two fixes offered, tagging entries holds up better than clearing the list on a non-re-exec start: clearing depends on startup order (the list must be cleared before any child starts, and the daemon would have to tell a re-exec start from a fresh one), while a tag is checked on every scan and does not depend on order. A pid alone is a weak tag because the new daemon could get the old daemon's pid after a reboot or after pids wrap. Adding the process start time from /proc/self/stat, which syscall.Exec keeps, tells the two apart, and it costs nothing more on Linux since the reaper already reads /proc. Off Linux the reaper still does nothing. The rest of reap-own-session-children stands: no wait4(-1), no picking drivers by session, goroutine Wait for the running image's drivers, ECHILD ignored, and a shared lock.

## statement

Each driver-list entry records the daemon that wrote it, not just the driver's pid: the driver's pid plus the recording daemon's pid and that daemon's process start time (field 22 of /proc/self/stat). A re-exec (syscall.Exec in place) keeps both, so drivers started before a re-exec stay the running daemon's; a fresh start after a crash, stop/start or systemd restart has a new pid or a new start time. The reaper, at start and on each tick and under the list lock, only calls wait4(pid, WNOHANG) on entries whose recorded daemon pid and start time match the running daemon's and whose pid is still its child (ppid from /proc). It drops every entry recorded by another daemon without waiting on it, since init has reaped or will reap those drivers. A stale pid that has since been reused by a git, make or hook child of the new daemon is therefore never waited on, whenever the scan runs, so the startup scan does not need to run before the daemon starts any child. The leaf's test adds the stale case: an entry recorded under another daemon's pid/start time whose pid is a live, unwaited child of the test process is dropped from the list and not reaped, and that child's own Wait still returns its exit status.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6q3wz-r3-f1](../question/review-01m6q3wz-r3-f1.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
