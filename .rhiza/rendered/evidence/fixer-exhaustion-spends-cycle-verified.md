# fixer-exhaustion-spends-cycle-verified

evidence · `01M657VXEMTFRSTBPEZPHMJSPK` · actor [claude-code-43542eb3](../actor/claude-code-43542eb3.md)

## result

PASS

## statement

At ca9f891, go test ./internal/run -run TestFixChangeSet -count=1 -v passes, including the new TestFixChangeSetExhaustionSpendsCycle; go test ./internal/run -count=3, make fmt-check and make vet also pass. The test has the change-set fixer of cycle 2 fail all 1 + MaxFixCycles (3) attempts. It asserts that alpha's cycle 2 outcome is CheckFailed with a detail starting 'cycle 2: fix failed: ', that no check session starts in cycle 2 for alpha, beta or gamma, that the run then reaches cycle 3, the last one under MaxFixCycles=2, with fixer sessions (2,1),(2,2),(2,3),(3,1), and that the step passes there. So the spent cycle counts against the cap. I read nextCheckActions in internal/run/checks.go (lines 214-228): when the fix phase fails it ends each unjudged check of the cycle with 'cycle N: fix failed: <reason>' before any checkStart, and changeSetFix retries before giving up. make self-check fails only on stale .rhiza/rendered files (the new actor's page and index.md), which is model render state and not this leaf's code.

## Links

### Out

- verifies [fixer-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
