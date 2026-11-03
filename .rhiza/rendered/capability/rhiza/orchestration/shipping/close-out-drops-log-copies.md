# close-out-drops-log-copies

capability · `01M6GD438PP0YPH7YYYZT2HAT4` · parent [shipping](../shipping.md)

## acceptance

A testscript ships a change set opened in the main checkout (run start --new or change new) and merges its PR. ship close-out exits 0 with main_updated true, the main checkout's HEAD is origin/main, and its log entries are tracked. A second case gives the main checkout an untracked .rhiza/log file whose bytes differ from the merged one: close-out exits 1, update-main is failed, and the file is unchanged

## statement

Before ship close-out fast-forwards the main checkout, it removes the untracked .rhiza/log files there that the incoming origin/main holds byte for byte, so a merged change set whose first entries were opened in the main checkout updates main. Untracked files that differ or that main lacks stay, and update-main fails and names them as it does now

## Links

### Out

- decided-by [close-out-untracked-log-and-retry](../../../../decision/close-out-untracked-log-and-retry.md)

### In

- verifies [close-out-drops-log-copies-verified](../../../../evidence/close-out-drops-log-copies-verified.md)
- verifies [close-out-drops-log-copies-qa-1](../../../../evidence/close-out-drops-log-copies-qa-1.md)
- verifies [pr-211-fix](../../../../evidence/pr-211-fix.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
