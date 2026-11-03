# close-out-untracked-log-and-retry

decision · `01M6GD0DBQCP1Q8S5K5VKV7Z4V`

## reasoning

Intent: m24-on-the-vm is the landed intent the previous Fix (folder trust, 01M6FSQ) serves, and this failure also hits the VM's unattended daemon. Approach: the main checkout must keep drafts' and accepted change sets' first entries, since the daemon reads in-flight change sets there, so writing them only in a worktree would hide them; dropping only files identical to incoming main loses nothing (main has the same bytes) and also fixes ship-change close-outs. Retry: internal/api/ship_close.go removes t.statePath() under r.Forget before shipMainBlocked and git pull --ff-only, so a failed pull leaves no state for a retry; moving the removal after step update-main makes close-out idempotent. All three follow the question's recommendations, which the code and the change set's intent support.

## statement

Change set 01M6GCP serves m24-on-the-vm. Before pulling main, close-out removes untracked .rhiza/log files in the main checkout that the incoming main has byte for byte, then pulls; untracked files that differ or are absent from main stay, and the pull reports them as it does now. `ship close-out --forget` deletes the ship state only after main is updated, so a close-out whose pull failed can be retried and finishes.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-63734d1c-1](../question/question-63734d1c-1.md)
- derived-from [close-out-fix-plan](../source/close-out-fix-plan.md)

### In

- decided-by [close-out-drops-log-copies](../capability/rhiza/orchestration/shipping/close-out-drops-log-copies.md)
- decided-by [close-out-retry-after-forget](../capability/rhiza/orchestration/shipping/close-out-retry-after-forget.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
