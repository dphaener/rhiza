# daemon-prunes-orphan-leftovers

capability · `01M7642M03G0506XQ7B2G4S0VF` · parent [daemon](../daemon.md)

## acceptance

Unit tests in the package that implements the sweep build fake index and scratch trees with live and orphaned entries, inject the clock, and show only the orphans removed and the freed bytes logged; go test -trimpath on that package.

## statement

The daemon prunes orphaned leftovers at most once an hour: index databases whose worktree path no longer exists, and agent-session scratch folders of sessions that belong to no live run and are older than a day; it never touches a live run's files and logs what it freed.

## Links

### In

- verifies [daemon-prunes-orphan-leftovers-verified](../../../../evidence/daemon-prunes-orphan-leftovers-verified.md)
- verifies [daemon-prunes-orphan-leftovers-qa-1](../../../../evidence/daemon-prunes-orphan-leftovers-qa-1.md)
- verifies [daemon-prunes-orphan-leftovers-verified-2](../../../../evidence/daemon-prunes-orphan-leftovers-verified-2.md)
- verifies [pr-467-runs-clean-up-rhiza-s-own](../../../../evidence/pr-467-runs-clean-up-rhiza-s-own.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
