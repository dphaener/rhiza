# reap-own-session-children

decision · `01M6Q4344CFT9EQRX2FK909RX3`

## reasoning

cmd/rhiza/detach_unix.go already reaps each driver with go cmd.Wait() (since #210). So the 50 zombies come from the daemon re-exec'ing itself in place (internal/daemon/reexec_unix.go, syscall.Exec): the waiting goroutines die and the new image does not know the children it had. The intent's wait4(-1, WNOHANG) would cover those, but the question is right that it would also reap git/make/hook children that os/exec goroutines (some from web UI requests) are waiting on, and their Wait would fail with 'no child processes'. That is a new bug in normal operation. Telling drivers apart by getsid(pid) == pid is not enough: internal/api/daemon_hook.go runs the owner's hooks through sh -c with no SysProcAttr, so a hook like 'setsid foo' makes the daemon's direct child a session leader (plan review round 1, f1). A pid list on disk names exactly the children the daemon detached, survives a re-exec because it lives in the git common dir, and a recorded pid that is still a zombie cannot be reused, so waiting on listed pids only never reaps a child os/exec waits on. Running it at startup plus on a timer means a re-exec'd daemon clears its leftovers right away. The daemon runs on the Linux VM, so reading /proc is enough.

## statement

The daemon never calls wait4(-1), and it does not pick drivers by session leadership: a hook or other os/exec child that runs setsid(1) calls setsid() in place and leads its own session too. Instead each detached driver's pid is recorded: detachProgram, after cmd.Start, appends the pid to a driver list under <git common dir>/rhiza/ (beside the run logs), which outlives a re-exec of the daemon. It keeps its go cmd.Wait() for drivers started by the running image, and drops the pid from the list once that Wait returns. A reaper runs once when the daemon starts (so after every re-exec) and then on a timer: for each pid on the list that is still a child of the daemon (on Linux from /proc: ppid equal to the daemon's pid) it calls wait4(pid, WNOHANG), and drops from the list each pid it reaps or that is no longer its child. Children run through os/exec for git, make and hooks are never on the list, so their cmd.Wait/Output keeps working, whatever they do with sessions. ECHILD from a reaper/goroutine race on the same driver is ignored. The one limit left is pid reuse in the instant between a goroutine's Wait reaping a driver and it dropping the pid from the list, which would need a new direct child of the daemon to draw that very pid; list writes and the reaper's scan share a lock to keep that window to the Wait itself. Off Linux the reaper does nothing, like detach_other. The test starts a Setsid child recorded on the list that is not waited on, lets it exit, runs the reaper, and checks the pid is gone (no zombie) and off the list, and that a plain os/exec child that calls setsid() itself, running at the same time, still gets its own exit status from Wait.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-ab537112-1](../question/question-ab537112-1.md)
- derived-from [daemon-reaper-plan](../source/daemon-reaper-plan.md)

### In

- decided-by [daemon-reaps-drivers](../capability/rhiza/orchestration/daemon/daemon-reaps-drivers.md)

## History

- `01M6Q3X16PDW20R7N57BJ5JK6G` Fix The daemon never reaps its finished run drivers: 50 zombie 'rhiza' processes, all children of the daemon (pid 1348, up 30 h), the oldest 30 h old; the owner saw '50 zombie processes' at ssh login on 2026-11-06. The daemon starts each run's driver (run start/resume --foreground) and never waits on it; it also re-execs itself in place on a new binary (same pid), after which it cannot know the children it had before. Reap with a periodic non-blocking wait for any child (wait4(-1, WNOHANG) in a loop, which also covers children from before a re-exec), and add a test that a finished child of the daemon leaves no zombie.
