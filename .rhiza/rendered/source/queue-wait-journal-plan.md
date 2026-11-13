# queue-wait-journal-plan

source · `01M79DDV1SRFJE2R3VCGSCP391`

## quotes

```json
[]
```

## raw_sha256

040c32c4b1e094f9453d120692b60eb7d0cff27d30494e0efd6422f902b3df14

## session

6785417f-0423-45fb-8791-01a1d2d346d4

## summary

Autonomous plan run for change set 01M76XSY (Queue wait journal), serving m35-ship-40m. No owner took part. The planner read the daemon's pick code: every poll works out a Skip reason for each accepted change set (internal/daemon/pick.go, daemonPoll1 in internal/api/daemon.go), but daemon.state.json keeps only the last poll's skips, so a wait leaves no history. The machine hold already reports as the Skip reasons load-high and memory-low (internal/daemon/load.go). The planner raised three questions: where the wait is journaled, how the perf analysis reads it, and whether it is trimmed. The arbiter answered with Decision queue-wait-journal-file: an append-only queue.wait.jsonl in the daemon's per-clone folder, a line only when a change set's wait starts, changes reason, or its run starts, a typed line and reader in internal/daemon, read directly by the perf analysis, no CLI command, never trimmed. The planner wrote two leaves under daemon-picks: the daemon's journal lines, and the line type and reader.

## title

Queue wait journal plan (autonomous)

## Links

### In

- derived-from [queue-wait-journal-file](../decision/queue-wait-journal-file.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
