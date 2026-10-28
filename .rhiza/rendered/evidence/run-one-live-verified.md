# run-one-live-verified

evidence · `01M62413SHBF6A0XR9BKNKHAT2` · actor [claude-code-5fb52618](../actor/claude-code-5fb52618.md)

## result

PASS

## statement

At a5204c9 I ran go test ./cmd/rhiza -run TestScripts/(run-one-live|run-status|review-signoff|autonomous-gate-rule) and go test -race -count=5 ./internal/api -run TestConcurrentStartsOneBegins: both pass. run-one-live.txtar stops a hold run with --no-wait, a second start (and its --dry-run) exits 2 with 'already has a live run, <id>' and 'next: rhiza run resume <id>', leaving one journal; after the first finishes a new start succeeds; with a failed flaky run and a second run stopped at its wait, run resume --retry (and --dry-run) on the failed run exits 2 naming the second. TestConcurrentStartsOneBegins races 8 starts on one change set: exactly one begins, the rest exit 2. Removing the refuseLiveRun check from startRun in a scratch worktree made both the Go test and the testscript fail. Read the diff: start and resume check LiveOf (unfinished journal, interrupted included) under a per-change-set flock, and the engine's Admit hook holds it while the run is reopened. make check exits 0.

## Links

### Out

- verifies [run-one-live](../capability/rhiza/workflows/workflow-runs/run-one-live.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
