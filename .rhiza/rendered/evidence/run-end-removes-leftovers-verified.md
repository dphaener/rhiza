# run-end-removes-leftovers-verified

evidence · `01M772T0N5T7XEFBRZYJ04V3Q9` · actor [claude-code-e758c10c](../actor/claude-code-e758c10c.md)

## result

PASS

## statement

At ee592a1c6 (attempt 2) I ran `go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/run-end-removes-leftovers` and it passed. The new script uses a fake rex and a fake claude. It starts a live run that stops at a wait step and an ended run that succeeds through the engine. The agent of each run makes a session scratch folder under $TMPDIR/claude-*/<cwd>/<sid> and an index database (via `rhiza slice`). Once the ended run finishes, its scratch folder is gone, and the live run's scratch folder and index stay, as does an unrelated folder. `rhiza run abandon` on the live run then removes its worktree, scratch folder and index, while the other worktree's index stays. This answers attempt 1's finding, which was that no testscript existed. The unit tests TestRemoveRunScratchTakesOnlyTheEndedRunsSessions, TestRemoveWorktreeIndexKeepsOtherIndexes and TestRemoveCheckoutRemovesItsIndex pass (internal/api). So do the internal/index, internal/adapter/... and internal/leftovers tests, and the run-abandon and ship-close scripts. I read the diff: the hooks are in runEngineRetry (st.Ended), runRunAbandon/abandonCleanup, runShipCloseOut and RemoveCheckout. The index removal happens before the worktree goes, because the database is named by the resolved path. Note: one broader internal/api run tripped the root-state guard (~/code/.rhiza-root/usage.live.json changed). Its tests had passed, and a rerun of the same Session tests was clean. This change only adds ScratchGlobs stubs to those tests, so the write came from the live root, not from this change.

## Links

### Out

- verifies [run-end-removes-leftovers](../capability/rhiza/workflows/workflow-runs/run-end-removes-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
