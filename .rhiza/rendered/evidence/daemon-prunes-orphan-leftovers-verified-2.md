# daemon-prunes-orphan-leftovers-verified-2

evidence · `01M775CAPFHE8NJ3H82HA9V3ZT` · actor [claude-code-9f53426f](../actor/claude-code-9f53426f.md)

## result

PASS

## statement

Verifier (attempt 2) at 27033141f: read the diff (daemon.go passes every served project's git common dir, collected after all registered projects are initialised in the poll, to daemonSweepLeftovers, which dedups and sweeps each project's session scratch; index sweep global). Ran go test -trimpath -count=1 ./internal/leftovers/... ./internal/index/... (ok) and ./internal/api -run TestDaemonSweepLeftovers -v: PASS. The test builds fake index DBs (live and orphaned worktree) and session scratch folders across two projects with an injected clock: only orphans go (ended-old, no-run-old, unknown-run-old, other project's crashed-run-old), live-run, fresh and other project's live-run scratch stay, stderr logs the freed bytes with 1 index DB and 4 folders, a sweep 59 min later does nothing, one at 60 min runs. A broader ./internal/api -run 'Leftover|Sweep|Daemon' run passed its tests but its root-state guard flagged ~/code/.rhiza-root/usage.live.json changing; that file holds a real usage capture written live by the VM's sessions at that moment, not by this change.

## Links

### Out

- verifies [daemon-prunes-orphan-leftovers](../capability/rhiza/orchestration/daemon/daemon-prunes-orphan-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
