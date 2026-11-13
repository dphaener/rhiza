# poll-behind-names-steps-verified

evidence · `01M7A1XD2N786832PFHKX6YF5X` · actor [claude-code-6150daeb](../actor/claude-code-6150daeb.md)

## result

PASS

## statement

At e2f56e49a re-ran the Checks: block: go test -trimpath ./internal/api -run TestWatchdogPollBehind (exit 0; TestWatchdogPollBehind and TestWatchdogPollBehindSteps PASS), go test -trimpath ./internal/api -run 'TestWatchdogPollBehind|Daemon|Poll' (exit 0, 259s), rhiza ship generate (exit 0), rhiza render (exit 0, tree unchanged); gofmt -l and go vet ./internal/api clean. Diff read: pollSteps times registry load, open, and per project reap/settings/machine/housekeeping/usage/demand/starts/notify/watchdog/events/model steps; projects poll sequentially so no race. TestWatchdogPollBehindSteps uses a fake clock advanced inside named steps and asserts the exact line keeps the 'watchdog: poll took ... poll interval' prefix and ends '; slowest: rhiza watchdog 5m58s, rhiza demand 21s, other events 3s' (slowest first, project and rounded time), that the recorded watchdog action carries it, and that a poll within the limit logs nothing. Only the four named Go files plus the builder's actor record and rendered files changed; no TODOs or stubs.

## Links

### Out

- verifies [poll-behind-names-steps](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-behind-names-steps.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
