# scratch-sweep-one-listing-verified

evidence · `01M7A1ST065DKCVWH1FZP6SSH9` · actor [claude-code-fe660995](../actor/claude-code-fe660995.md)

## result

PASS

## statement

At 6defe991e re-ran the commit's Checks: go test -trimpath ./internal/leftovers/... ./internal/adapter/... (exit 0); go test -trimpath ./internal/api -run 'TestDaemonSweepLeftovers|TestRemoveRunScratch|TestSessionInterview|TestSessionHook' (exit 0: TestDaemonSweepLeftoversPrunesOnlyOrphans over the interview adapter's $HOME/scratch/<sid> layout, TestRemoveRunScratchTakesOnlyTheEndedRunsSessions, TestSessionInterview*, TestSessionHook* pass); rhiza ship generate (exit 0, 0 files changed); rhiza render (exit 0, no tracked diff). Read the diff: Adapter.ScratchGlobs is replaced by ScratchParents(getenv) in adapter.go, claudecode (<tmp>/claude-*/*), the interview fake ($HOME/scratch) and nil for ask/adapter_test fakes; leftovers.Scratch groups orphaned safe IDs by agent, globs each parent once, lists each through an injectable ReadDir and removes only matching child folders older than ScratchMaxAge, counting bytes; run_leftovers.go removes <parent>/<id> via leftovers.RemoveSession, which refuses unsafe IDs. TestScratchReadsEachParentOnce records 3000 sessions (1500 folders in 3 parents, 10% of live runs, plus 1500 ghosts), counts reads through a wrapper, asserts each of the 3 directories is read exactly once and 1350 folders / 6750 bytes freed with the live run's 50 left; TestScratchRemovesOnlyOldOrphans keeps fresh and live folders; TestRemoveSessionSafeIDOnly covers ../ and * IDs; TestScratchParents covers the Claude Code parents. Only the leaf's listed files plus the actor record and rendered files changed; no TODOs or stubs.

## Links

### Out

- verifies [scratch-sweep-one-listing](../capability/rhiza/orchestration/daemon/scratch-sweep-one-listing.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
