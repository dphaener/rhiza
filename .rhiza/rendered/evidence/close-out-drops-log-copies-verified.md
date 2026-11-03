# close-out-drops-log-copies-verified

evidence · `01M6GSMQNNQB35HKZQS2V6GY82` · actor [claude-code-093c7fed](../actor/claude-code-093c7fed.md)

## result

PASS

## statement

At 7cf32ec6 I read the diff of cfe2054e and 7cf32ec6. Before the ff-only pull, dropIncomingLogCopies in internal/api/ship_close.go fetches origin/main, lists the untracked .rhiza/log files and removes only those whose hash-object equals origin/main:<path>. I ran go test ./cmd/rhiza -run TestScripts/ship-close-out -count=1 and it passed. In ship-close-out-log-copies.txtar, case 1 opens a change set with change new in the main checkout and merges it through a clone. close-out exits 0 with main_updated true, HEAD equals the pushed origin/main, and git status shows no untracked .rhiza/log. Case 2 gives the main checkout differing untracked log files. close-out exits non-zero with main_updated false, step update-main and result failed. The untracked list is unchanged, and cmp shows each file's bytes are the same. This answers attempt 1's finding. Mutation check: in a scratch copy with the dropIncomingLogCopies call removed, the script fails at case 1 (main_updated false). In make check, fmt-check, vet and test pass. self-check fails only on a stale .rhiza/rendered/actor/claude-code-f7650bc4.md and index.md. That file came from 3d6fd44c on the other leaf, close-out-retry-after-forget, not from this leaf's diff.

## Links

### Out

- verifies [close-out-drops-log-copies](../capability/rhiza/orchestration/shipping/close-out-drops-log-copies.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
