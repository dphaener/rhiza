# poll-overrun-logged-live-qa-3

evidence · `01M7AAFANNSS0E72ZTYYJWVH68` · actor [claude-code-10e530e4](../actor/claude-code-10e530e4.md)

## result

PASS

## statement

QA cycle 3 at 1820b5bd5: built rhiza to a scratch path; hermetic root with project repo and a git wrapper delaying the first poll's settings read by 20s (value kept). rhiza daemon --once: watchdog on (default) logged one 'poll still running after 15s, in repo settings' then the poll-took line; rhiza.daemon.watchdog.enabled=false in git config or in rhiza.toml (rhiza daemon set) logged nothing; pollBehind=6 (30s limit) logged nothing; pollBehind=1 logged 'after 5s, in repo settings' once. Long-running daemon with the hang in the second poll: exactly one in-flight line, then the poll-took line. go test ./internal/api -run 'TestWatchdogPoll|TestWatchdogPolicy' and ./cmd/rhiza TestScripts/watchdog-policy pass

## Links

### Out

- verifies [poll-overrun-logged-live](../capability/rhiza/orchestration/daemon/daemon-watchdog/poll-overrun-logged-live.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
