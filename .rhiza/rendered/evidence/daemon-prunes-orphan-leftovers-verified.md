# daemon-prunes-orphan-leftovers-verified

evidence · `01M774FEVX8J3DFJAW5HCZ48XB` · actor [claude-code-7de91295](../actor/claude-code-7de91295.md)

## result

PASS

## statement

Verifier at 96ab8695b read the leaf's diff (633298a41, 96ab8695b): index.Open records the resolved worktree root in meta; index.PruneOrphans removes .db (+wal/shm) whose recorded root is gone and keeps live and unrecorded/unreadable ones; leftovers.Scratch removes session scratch folders older than 24h whose run is not live (skips live-run sessions and fresh folders) and sums bytes; daemonSweepLeftovers runs from the daemon's home step on a SweepSchedule (at start, then SweepEvery=1h) and logs one stderr line with freed bytes only when something went. Ran go test -trimpath -count=1 ./internal/leftovers/... ./internal/index/... ./internal/daemon/... (all ok) and go test -trimpath -count=1 ./internal/api -run TestDaemonSweepLeftovers -v: PASS. That test builds fake index DBs (live, gone) and scratch folders (live-run old, ended old, ended fresh, no run, unknown run) with an injected clock: only the 3 orphans and the gone index are removed, stderr carries the exact freed bytes, a sweep 59 minutes later does nothing and one at 60 minutes sweeps again, never touching the live run's folder. A wider api run (-run Index|Leftover) passed every test but rootguard flagged .rhiza-root/usage.live.json changed; that file is rewritten each second by live Claude Code status lines on the host and is not in rootguard's owned list, unrelated to this leaf.

## Links

### Out

- verifies [daemon-prunes-orphan-leftovers](../capability/rhiza/orchestration/daemon/daemon-prunes-orphan-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
