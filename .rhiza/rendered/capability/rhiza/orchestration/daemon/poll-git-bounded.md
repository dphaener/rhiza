# poll-git-bounded

capability · `01M7A01HZMHFKMQ821D1SR6BFH` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/api -run 'TestPollGit': with a fake ctx.Start whose git never exits and a fake clock or a shortened limit, a test runs a project context's Git and GitConfig and asserts each returns an error naming the command and the limit after the limit, the program was killed, and a git that exits in time returns its output unchanged; a second test calls daemonHome with ctx.Root at the home project's own dir (so it returns that project's ctx, not a projectCtx) and the same fake ctx.Start, and asserts the home ctx's Git and GitConfig the poll uses are bounded the same way, not the caller's unbounded Git. Touches internal/api/daemon_projects.go, internal/api/daemon.go, internal/api/daemon_projects_test.go

## statement

Every git program the daemon's poll runs through a project's context (the home project's Git and each other project's Git and GitConfig) is killed after 30 s with GIT_TERMINAL_PROMPT=0, as execWithin does for the model-branch work, and fails with an error naming the git command and the time limit, so a hung git costs a poll at most that long per call; the self-update's make and the run drivers the poll detaches are not bounded

## Links

### Out

- decided-by [decision-poll-git-timeout](../../../../decision/decision-poll-git-timeout.md)

### In

- verifies [poll-git-bounded-verified](../../../../evidence/poll-git-bounded-verified.md)
- verifies [poll-git-bounded-qa-1](../../../../evidence/poll-git-bounded-qa-1.md)
- verifies [poll-git-bounded-qa-2](../../../../evidence/poll-git-bounded-qa-2.md)
- verifies [poll-git-bounded-qa-3](../../../../evidence/poll-git-bounded-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
