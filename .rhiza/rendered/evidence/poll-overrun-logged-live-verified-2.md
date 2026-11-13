# poll-overrun-logged-live-verified-2

evidence · `01M7A8N4SGD7AAS1TBFZT6BHJP` · actor [claude-code-4bf6c1df](../actor/claude-code-4bf6c1df.md)

## result

PASS

## statement

Attempt 2 at 9ff80444b: acceptance `go test -trimpath ./internal/api -run TestWatchdogPollOverrun` exit 0 (covers TestWatchdogPollOverrun: one in-flight line naming step and project before release, none within limit, none with watchdog or poll-behind off; and TestWatchdogPollOverrunFirstPoll: re-arm from zero to enabled logs once at the new limit, re-arm to off or after stop logs nothing). Checks: `go test -trimpath -race -count=1 ./internal/api -run TestWatchdogPoll` exit 0. gofmt clean, go vet ./internal/api clean, `make self-check` exit 0 in a clean detached worktree with a clean env (the in-place failures were only this session's unrendered actor and untracked logs). Diff read: daemon starts with defaultWatchdogPolicy so the first poll is watched; atHome re-arms the poll's watch synchronously from prepare in the poll goroutine, timed from the poll's start, at most one line per poll; a stale rearm after stop is a no-op. watchdog_policy.go gains only the defaultWatchdogPolicy helper the fix needs. No TODOs or stubs; tests use fakes only.

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
