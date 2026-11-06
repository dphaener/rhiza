# test-scripts-sharded-verified

evidence · `01M6Q8XWSKGVJVCWJDSR1P8E15` · actor [claude-code-e8b74584](../actor/claude-code-e8b74584.md)

## result

PASS

## statement

At 662fcadb, go test -trimpath ./cmd/rhiza -run TestScriptShard -v -count=1 passed: for n=4 each shard's entries equal the round-robin picks over the sorted 271 testdata/script names, the shards are disjoint and cover all names, unset returns nil (testscript runs every script), and 0/4, 5/4, x, 1/0, 1/, /4, 1/4x, -1/4 are rejected. End to end: RHIZA_TEST_SHARD=x fails TestScripts with 'RHIZA_TEST_SHARD="x": want i/n with 1 <= i <= n'; with -run TestScripts/describe-golden$ (sorted index 65), shard 2/4 runs and passes it while 1/4, 3/4, 4/4 run nothing, and unset runs it. Diff read: Files passed with Dir cleared only when a shard is set; gofmt and go vet clean.

## Links

### Out

- verifies [test-scripts-sharded](../capability/rhiza/integrity/self-check/test-scripts-sharded.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
