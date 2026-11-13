# scratch-sweep-one-listing

capability · `01M7A01HP9C0SQC3RZMW0R2R18` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/leftovers/... ./internal/adapter/... and go test -trimpath ./internal/api -run 'TestDaemonSweepLeftovers|TestRemoveRunScratch|TestSessionInterview|TestSessionHook': a leftovers test records 3000 sessions (some of live runs) over a scratch tree of 1500 folders, counts directory reads through an injected reader or a counting wrapper, and asserts each scratch directory is read at most once while exactly the orphaned folders older than a day are removed and the freed bytes are counted; the claudecode test covers the scratch parents and session ID matching it now gives (no parents widened by an unsafe ID); TestDaemonSweepLeftoversPrunesOnlyOrphans still passes over the interview adapter's one-level $HOME/scratch/<sid> layout; a run's end still removes its sessions' scratch folders. Touches internal/leftovers/leftovers.go, internal/leftovers/leftovers_test.go, internal/adapter/adapter.go, internal/adapter/adapter_test.go, internal/adapter/claudecode/claudecode.go, internal/adapter/claudecode/claudecode_test.go, internal/api/daemon_leftovers.go, internal/api/daemon_leftovers_test.go, internal/api/run_leftovers.go, internal/api/session_interview_test.go, internal/api/session_ask_test.go

## statement

The leftovers sweep finds the session scratch folders to remove by listing each scratch level once per sweep and matching folder names against the recorded sessions of no live run, so its directory reads do not grow with the number of recorded sessions; what it removes and the line it logs are unchanged. The levels come from each adapter: the Adapter interface's ScratchGlobs gives way to a per-adapter list of scratch parents, glob patterns of the folders whose direct children are session folders named by session ID (Claude Code's <tmp>/claude-*/*, the interview test adapter's $HOME/scratch, none for the ask and adapter_test fakes), and the run-end cleanup in run_leftovers.go removes <parent>/<session ID> under the same parents

## Links

### Out

- decided-by [decision-scratch-sweep-one-listing](../../../../decision/decision-scratch-sweep-one-listing.md)

### In

- depends-on [leftovers-sweep-off-poll](leftovers-sweep-off-poll.md)
- verifies [scratch-sweep-one-listing-verified](../../../../evidence/scratch-sweep-one-listing-verified.md)
- verifies [scratch-sweep-one-listing-qa-1](../../../../evidence/scratch-sweep-one-listing-qa-1.md)
- verifies [scratch-sweep-one-listing-qa-2](../../../../evidence/scratch-sweep-one-listing-qa-2.md)
- verifies [scratch-sweep-one-listing-qa-3](../../../../evidence/scratch-sweep-one-listing-qa-3.md)
- verifies [pr-543-fix-the-root-daemon-s-poll](../../../../evidence/pr-543-fix-the-root-daemon-s-poll.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
