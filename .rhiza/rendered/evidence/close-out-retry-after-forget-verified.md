# close-out-retry-after-forget-verified

evidence · `01M6GZ8AJ3FZ5HTKZA7VC5WDPS` · actor [claude-code-572a7a96](../actor/claude-code-572a7a96.md)

## result

PASS

## statement

At 8f2547e4, cmd/rhiza/testdata/script/ship-close-out-forget-retry.txtar dirties a tracked file in main and runs ship close-out --forget: it exits non-zero with main_updated false and update-main failed, and .git/rhiza/ship/<cs>.json still exists. After git checkout -- file.txt the same command reports main_updated true and the state file is gone. go test ./cmd/rhiza -run TestScripts/ship-close-out-forget-retry -count=1 passes. A manual scratch repro with a fresh build gave exit=1 (update-main failed, state file kept), then exit=0 (main_updated true, state file removed). The diff in internal/api/ship_close.go moves the os.Remove of the state path to after update-main succeeds. Attempt 1's finding is fixed: fmt-check, go vet ./... and make self-check (rhiza check) all pass with RHIZA_* unset, and go test ./cmd/rhiza -count=1 passes. The only failures under heavy machine load (load average ~15 on 4 cores) were timing budgets (TestReadBudget, slice TestBudget, modelapp) and command timeouts, which report contention from other packages and do not touch this diff.

## Links

### Out

- verifies [close-out-retry-after-forget](../capability/rhiza/orchestration/shipping/close-out-retry-after-forget.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
