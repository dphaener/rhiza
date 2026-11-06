# test-scripts-sharded

capability · `01M6Q5HHTX26CPP16DNNEK7T80` · parent [self-check](../self-check.md)

## acceptance

A Go test in cmd/rhiza (go test -trimpath ./cmd/rhiza -run TestScriptShard) asserts that for n=4 the four shards are disjoint, together cover every script in testdata/script, and are round-robin over the sorted names; that unset selects all scripts; and that a malformed value (0/4, 5/4, x) fails with a clear error.

## statement

TestScripts runs one shard of cmd/rhiza's scripts when RHIZA_TEST_SHARD=i/n is set (1-based): the sorted testdata/script/*.txtar names whose index mod n is i-1, passed as testscript Files; unset, it runs every script as before.

## Links

### Out

- decided-by [ci-speed-shape](../../../../decision/ci-speed-shape.md)

### In

- depends-on [make-test-shard-targets](make-test-shard-targets.md)
- verifies [test-scripts-sharded-verified](../../../../evidence/test-scripts-sharded-verified.md)
- verifies [test-scripts-sharded-qa-1](../../../../evidence/test-scripts-sharded-qa-1.md)
- verifies [pr-322-ci-takes-11-18-minutes-per-pr](../../../../evidence/pr-322-ci-takes-11-18-minutes-per-pr.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
