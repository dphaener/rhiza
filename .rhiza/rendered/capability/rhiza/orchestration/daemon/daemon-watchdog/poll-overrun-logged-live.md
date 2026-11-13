# poll-overrun-logged-live

capability · `01M7A01J6NM4S5VW0XRCPK60ZP` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogPollOverrun: a test runs a poll whose step blocks on a channel past a shortened limit and asserts exactly one in-flight line naming that step and project is written before the step is released, none for a poll within the limit, and none with poll-behind disabled. Touches internal/api/watchdog_pollbehind.go, internal/api/watchdog_pollbehind_test.go, internal/api/daemon.go

## statement

While a poll is still running past the poll-behind limit, the daemon writes one log line naming the step and project it is in and how long the poll has run, as in 'watchdog: poll still running after 15s, in rhiza watchdog', so the log shows a stuck step before the poll ends; it writes at most one such line per poll and none when the watchdog or its poll-behind check is off

## Links

### Out

- depends-on [poll-behind-names-steps](poll-behind-names-steps.md)
- decided-by [decision-slow-step-in-log](../../../../../decision/decision-slow-step-in-log.md)

### In

- verifies [poll-overrun-logged-live-verified](../../../../../evidence/poll-overrun-logged-live-verified.md)
- verifies [poll-overrun-logged-live-qa-1](../../../../../evidence/poll-overrun-logged-live-qa-1.md)
- verifies [poll-overrun-logged-live-verified-2](../../../../../evidence/poll-overrun-logged-live-verified-2.md)
- verifies [poll-overrun-logged-live-qa-2](../../../../../evidence/poll-overrun-logged-live-qa-2.md)
- verifies [poll-overrun-logged-live-verified-5](../../../../../evidence/poll-overrun-logged-live-verified-5.md)
- verifies [poll-overrun-logged-live-qa-3](../../../../../evidence/poll-overrun-logged-live-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
