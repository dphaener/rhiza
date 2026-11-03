# close-out-retry-after-forget

capability · `01M6GD4DWDYX2T29SPWWRY2DF6` · parent [shipping](../shipping.md)

## acceptance

A testscript runs ship close-out --forget while the main checkout has uncommitted tracked changes: it exits 1, update-main is failed, and .git/rhiza/ship/<cs>.json still exists. After the changes are discarded, the same close-out --forget exits 0 with main_updated true and the ship state file is gone

## statement

ship close-out --forget deletes the change set's ship state only after update-main succeeds, so a close-out whose pull failed keeps its state and finishes when retried

## Links

### Out

- decided-by [close-out-untracked-log-and-retry](../../../../decision/close-out-untracked-log-and-retry.md)

### In

- verifies [close-out-retry-after-forget-verified](../../../../evidence/close-out-retry-after-forget-verified.md)
- verifies [close-out-retry-after-forget-qa-1](../../../../evidence/close-out-retry-after-forget-qa-1.md)
- verifies [pr-211-fix](../../../../evidence/pr-211-fix.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
