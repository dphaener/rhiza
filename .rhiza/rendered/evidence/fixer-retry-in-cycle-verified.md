# fixer-retry-in-cycle-verified

evidence · `01M657TEVE3QAHKRN9MH8RWKXG` · actor [claude-code-325e86f6](../actor/claude-code-325e86f6.md)

## result

PASS

## statement

At e59642b I read the diff: changeSetFix (internal/run/checks.go) restarts a change-set fixer judged fail within the same cycle at attempt+1, up to 1+MaxFixCycles attempts, before the cycle counts as spent. Attempt is journalled on review-fix and review-fix-judged entries. go test ./internal/run -run TestFixChangeSet -count=1 passed, including the new TestFixChangeSetRetriedInCycle: a fixer fails attempt 1 and passes attempt 2 in cycle 2, the run ends review:pass with 2 cycles (no new cycle), and the review-fix-judged pass entry comes before cycle 2's check-session. With the retry condition disabled in a scratch copy (git archive), that test failed (outcomes = [review:fail]). go test ./... and make fmt-check vet test passed; one earlier make test run failed only in the timing-sensitive read benchmarks. make self-check failed only on stale .rhiza/rendered actor and index files from this run's sessions, not on this leaf's code.

## Links

### Out

- verifies [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
