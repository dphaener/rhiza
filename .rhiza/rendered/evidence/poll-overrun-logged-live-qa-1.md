# poll-overrun-logged-live-qa-1

evidence · `01M7A5SQWATHPGXWH0S0RCHZT6` · actor [claude-code-51810022](../actor/claude-code-51810022.md)

## result

FAIL

## statement

Built binary with a git hung in the second poll's settings step: one 'watchdog: poll still running after 15s, in repo settings' line before the poll-took line; none with rhiza.daemon.watchdog.enabled=false. But the same hang in the first poll of a daemon image (default watchdog on, pollBehind 3) logged only the poll-took line after 30s and no in-flight line: the watcher is armed with homePol before that poll has read the settings.

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
