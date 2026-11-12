# run-end-removes-leftovers-qa-1

evidence · `01M774R4MB2EZ2DMYF3AM3P3Y5` · actor [claude-code-98a70447](../actor/claude-code-98a70447.md)

## result

PASS

## statement

QA at c90d32c99 with a binary built from the branch (go build -trimpath), run by hand in a scratch Rhiza root with its own HOME, XDG_CACHE_HOME and TMPDIR, plus the script's fake rex and fake claude (the fake makes $TMPDIR/claude-test/<cwd>/<session-id>/scratchpad and runs rhiza slice so the worktree gets an index). I started a build-hold run (stopped at its wait, --no-wait) and a build-done run that succeeded. After the success, the ended run's session scratch folder was gone; the live run's scratch folder and both index DBs stayed, as did an unrelated $TMPDIR/claude-test/other/keep. `rhiza run abandon` on the live run printed 'remove-worktree done' and removed its worktree, its session scratch folder and its index DB, while the ended run's index stayed. I also checked that the real Claude Code layout on this host ($TMPDIR=/home/rhiza/tmp, /home/rhiza/tmp/claude-1000/<encoded cwd>/<session id>) is matched by the adapter's claude-*/*/<id> pattern.

## Links

### Out

- verifies [run-end-removes-leftovers](../capability/rhiza/workflows/workflow-runs/run-end-removes-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
