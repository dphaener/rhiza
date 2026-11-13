# poll-overrun-logged-live-qa-2

evidence · `01M7A9DFREC9T77YZQGF5A7VA0` · actor [claude-code-0fc16d12](../actor/claude-code-0fc16d12.md)

## result

FAIL

## statement

QA cycle 2: hung-git run wrote exactly one 'watchdog: poll still running after 15s, in p1 demand' while the poll ran; none with the watchdog off and the hang in demand. But with rhiza.daemon.watchdog.enabled=false and a git config read hung in the first poll's p1 settings step, the daemon still wrote 'watchdog: poll still running after 15s, in p1 settings' (first poll watched under shipped defaults although the startup already read the settings)

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
