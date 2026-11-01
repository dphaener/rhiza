# run-abandon

capability · `01M61R7GG0GGBSVKV2585MNPCW` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript abandons an interrupted fake run: run status shows it failed with the abandoned detail and the actor, its clean worktree and local branch are gone, rhiza run start on the change set then succeeds, rhiza run resume on the abandoned run exits 2 naming it, and abandoning a succeeded or already abandoned run exits 2; describe-golden is regenerated.

## statement

rhiza run abandon <run> journals an unfinished or failed run as failed with the detail 'abandoned by <actor>', removing its clean worktree and local branch unless another run still uses them (see abandon-clears-worktree), so a new run can start on the change set; it refuses a succeeded or already abandoned run, and rhiza run resume refuses an abandoned run, naming it.

## Links

### Out

- depends-on [run-one-live](run-one-live.md)
- decided-by [answer-8cd68bc7-12](../../../../decision/answer-8cd68bc7-12.md)

### In

- depends-on [abandon-failed-run](abandon-failed-run.md)
- verifies [run-abandon-verified](../../../../evidence/run-abandon-verified.md)
- verifies [run-abandon-qa-1](../../../../evidence/run-abandon-qa-1.md)
- verifies [pr-50-a-change-set-has-at-most-one](../../../../evidence/pr-50-a-change-set-has-at-most-one.md)
- verifies [pr-66-inbox-clears-stale-failures](../../../../evidence/pr-66-inbox-clears-stale-failures.md)
- verifies [pr-164-abandon-clears-worktree](../../../../evidence/pr-164-abandon-clears-worktree.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch
