# run-end-removes-leftovers

capability · `01M7642KZP2AFAVV6DMMMDK8JR` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript (and internal/api or internal/index unit tests) ends a run whose fake agent sessions made scratch folders and whose worktree has an index database, and shows both gone afterwards, while a live run's scratch folders and index stay.

## statement

When a run ends (done, failed or abandoned) Rhiza removes the scratch folders of the agent sessions it started for that run, and when a run's worktree is removed Rhiza removes that worktree path's index database under the user cache dir (rhiza/index).

## Links

### In

- verifies [run-end-removes-leftovers-verified](../../../../evidence/run-end-removes-leftovers-verified.md)
- verifies [run-end-removes-leftovers-qa-1](../../../../evidence/run-end-removes-leftovers-qa-1.md)
- verifies [pr-467-runs-clean-up-rhiza-s-own](../../../../evidence/pr-467-runs-clean-up-rhiza-s-own.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
