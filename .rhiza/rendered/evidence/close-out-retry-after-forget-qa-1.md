# close-out-retry-after-forget-qa-1

evidence · `01M6GZGX7DFRKH9M31JPAJKS3X` · actor [claude-code-4706da01](../actor/claude-code-4706da01.md)

## result

PASS

## statement

At bf5b09b2, with a scratch build in a scratch repo: ship close-out --forget failed at update-main (exit 1, main_updated false) because a differing untracked log file blocked the pull, and .git/rhiza/ship/<cs>.json was still there. I removed the blocking file and re-ran the same ship close-out --forget. It exited 0 with main_updated true, main HEAD equal to origin/main, and the state file gone. In the first case, a successful --forget close-out also removed its state file.

## Links

### Out

- verifies [close-out-retry-after-forget](../capability/rhiza/orchestration/shipping/close-out-retry-after-forget.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
