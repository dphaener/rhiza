# poll-overrun-logged-live-verified-5

evidence · `01M7AA2PBPD0PS6SAPT08XXBD7` · actor [claude-code-73bddaf5](../actor/claude-code-73bddaf5.md)

## result

PASS

## statement

Attempt 5 at c8a22036d. Acceptance `go test -trimpath ./internal/api -run TestWatchdogPollOverrun` exit 0 (TestWatchdogPollOverrun and TestWatchdogPollOverrunFirstPoll pass: one in-flight line naming 'in rhiza watchdog' and 'after 20ms' before the blocked step is released, none within the limit, none with the watchdog or poll-behind off). Commit Checks re-run: go vet ./internal/api exit 0; go test -trimpath -race -count=1 ./internal/api -run 'TestWatchdogPoll|TestWatchdogPolicy' exit 0; go test -trimpath ./cmd/rhiza -run TestScripts/watchdog-policy exit 0; rhiza ship generate and rhiza render exit 0 with no tracked diff. Daemon --once e2e reproduced with scratch builds of HEAD^ and HEAD, a hermetic root with rhiza.daemon.watchdog.enabled in rhiza.toml and a fake git hanging 20s on the second rhiza.project.watchdog.enabled read: base with watchdog off logged 'poll still running after 15s, in repo settings'; fixed with watchdog off logged nothing; fixed with watchdog on logged that line and then the poll-took line; all exit 0. Diff read: homePol is seeded from watchdogPolicyOf(conf.rows), the startup daemonConfOf(homeCtx) settings; watchPollOverrun logs at most once per poll and rearm keeps the limit timed from poll start. Leaf commits touch only internal/api daemon.go, watchdog_pollbehind.go, its test and watchdog_policy.go (the policy helper the first-poll fix needs), plus .rhiza records. No stubs or TODOs.

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
