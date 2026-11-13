# poll-git-bounded-verified

evidence · `01M7A4CX2SYQJPA0RRFN5N3E9H` · actor [claude-code-23679dfa](../actor/claude-code-23679dfa.md)

## result

PASS

## statement

At b76122f41 (with d2fd2bd84) re-ran the Checks block: go test -trimpath ./internal/api -run 'TestPollGit' exit 0 (TestPollGitBounded, TestPollGitBoundedHome pass); go test -trimpath ./internal/api -run 'TestDaemon|TestPollGit|TestFetchMain|TestExecWithin|TestModel' exit 0; rhiza ship generate exit 0; rhiza render exit 0 with no tracked change. Read the diff: projectCtx and daemonHome's own-project branch both route Git and GitConfig through boundedExecCtx (execWithin, GIT_TERMINAL_PROMPT=0) with pollGitTimeout 30 s; a timeout returns timeoutError naming the command and limit, and GitConfig now passes it through instead of 'not set'. Tests use a fake ctx.Start whose git hangs, a 10ms limit, assert both errors name 'git -C <dir> ... timed out after 10ms', every started git was killed and had GIT_TERMINAL_PROMPT=0, a fast fake git returns its output unchanged, and daemonHome at the home project's own dir yields a bounded ctx (fake unbounded Git/GitConfig/Exec t.Fatal if used). Changes confined to internal/api/daemon_projects.go, its test and the .rhiza log/rendered files; no TODOs or stubs; fakes only.

## Links

### Out

- verifies [poll-git-bounded](../capability/rhiza/orchestration/daemon/poll-git-bounded.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
