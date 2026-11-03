# close-out-drops-log-copies-qa-1

evidence · `01M6GZGFB67X0R8GRCXJ6HEMA6` · actor [claude-code-4706da01](../actor/claude-code-4706da01.md)

## result

PASS

## statement

At bf5b09b2 I built rhiza into a scratch dir and ran it by hand in a scratch repo with a bare origin, RHIZA_* unset. I opened a change set with change new in the main checkout, which left two untracked .rhiza/log dirs. I copied them into a git worktree on branch cs, added a feature commit, pushed it, and squash-merged it into origin/main from a second clone. ship close-out --forget exited 0 with main_updated true and update-main done. Main HEAD equalled origin/main, git status was clean, and the log entries were tracked. Second case: I appended bytes to one untracked log file that origin/main also held. close-out exited 1 with update-main failed, and the detail named that file. The file stayed untracked with its tampered bytes (cmp equal). The matching copies beside it were dropped. Third case: an untracked .rhiza/log/stray.json that origin lacks stayed put, and close-out still updated main. go test ./cmd/rhiza -run TestScripts/ship-close-out -count=1 passes.

## Links

### Out

- verifies [close-out-drops-log-copies](../capability/rhiza/orchestration/shipping/close-out-drops-log-copies.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
