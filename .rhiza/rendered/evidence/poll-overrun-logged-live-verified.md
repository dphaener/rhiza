# poll-overrun-logged-live-verified

evidence · `01M7A4HBAVWDXTHD9PYPPQVNMB` · actor [claude-code-e04b0978](../actor/claude-code-e04b0978.md)

## result

PASS

## statement

At 35e6f758: re-ran the commit's Checks. go test -trimpath ./internal/api -run TestWatchdogPollOverrun exits 0, also under -race -count=5 with TestWatchdogPoll*. rhiza ship generate and rhiza render both exit 0. go vet and gofmt are clean. The diff touches only watchdog_pollbehind.go, its test, daemon.go and .rhiza records. watchPollOverrun arms a time.AfterFunc at PollBehind x interval. When it fires it logs one 'watchdog: poll still running after <limit>, in <project> <step>' line, using pollSteps' mutex-guarded running stack. Stop suppresses the line, and a disabled watchdog or pollBehind<1 arms nothing. daemon.go defers it around each poll and writes the line to stderr. The test blocks a step on a channel and asserts exactly one line naming 'rhiza watchdog' and 20ms before release. It also asserts no line within the limit and none when the check is disabled (Enabled=false or PollBehind=0). No stubs or TODOs.

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
