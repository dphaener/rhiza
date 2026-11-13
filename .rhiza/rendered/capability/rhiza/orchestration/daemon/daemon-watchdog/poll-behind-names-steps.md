# poll-behind-names-steps

capability · `01M7A01J38C50QZG2X4JNEC3S8` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogPollBehind: with a fake clock that advances inside named steps, a test runs the poll-behind check and asserts the logged line and the recorded watchdog action name the slowest steps in order with their project and rounded time, that a poll within the limit logs nothing, and that the line keeps its existing prefix. Touches internal/api/watchdog_pollbehind.go, internal/api/watchdog_pollbehind_test.go, internal/api/daemon.go, internal/api/daemon_projects.go

## statement

The daemon times each step of a poll (registry load, then per project: settings, machine and usage reads, demand, starts, notify, watchdog, events and model-branch work, plus the home project's housekeeping), and when the poll-behind check fires, its log line and watchdog action name the slowest steps with their project and time, as in 'poll took 6m28s, over its limit of 15s (3 times the 5s poll interval); slowest: rhiza watchdog 5m58s, rhiza demand 21s'

## Links

### Out

- decided-by [decision-slow-step-in-log](../../../../../decision/decision-slow-step-in-log.md)

### In

- depends-on [poll-overrun-logged-live](poll-overrun-logged-live.md)
- verifies [poll-behind-names-steps-verified](../../../../../evidence/poll-behind-names-steps-verified.md)
- verifies [poll-behind-names-steps-qa-1](../../../../../evidence/poll-behind-names-steps-qa-1.md)
- verifies [poll-behind-names-steps-qa-2](../../../../../evidence/poll-behind-names-steps-qa-2.md)
- verifies [poll-behind-names-steps-qa-3](../../../../../evidence/poll-behind-names-steps-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
