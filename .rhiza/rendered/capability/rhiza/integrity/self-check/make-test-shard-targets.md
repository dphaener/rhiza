# make-test-shard-targets

capability · `01M6Q5HHVN71122BKP8ZG1DQPR` · parent [self-check](../self-check.md)

## acceptance

A Go test beside TestMakeTestTimeout in cmd/rhiza/makeupdate_test.go (go test -trimpath ./cmd/rhiza -run 'TestMake') runs make -n on each new target and asserts its go test line carries -trimpath and the budget skip or run regex, that test-api's shards for n=2 partition internal/api's listed tests with no budget test among them, and that test-budgets prints the timings; TestMakeTestTimeout still passes unchanged for make test.

## statement

The Makefile has one target per CI job: test-scripts SHARD=i/n (cmd/rhiza's scripts, budgets skipped), test-api SHARD=i/n (internal/api's top-level tests from go test -list, budgets excluded, split round-robin into an anchored -run regex), test-rest (every other package, budgets skipped), test-budgets (the timing budgets alone, -count=1 -p 1, printing the read-budget timings) and lint (fmt-check, vet, self-check); make test and make check still run the full suite.

## Links

### Out

- depends-on [test-scripts-sharded](test-scripts-sharded.md)
- decided-by [ci-speed-shape](../../../../decision/ci-speed-shape.md)

### In

- depends-on [ci-yml-parallel-jobs](ci-yml-parallel-jobs.md)
- verifies [make-test-shard-targets-verified](../../../../evidence/make-test-shard-targets-verified.md)
- verifies [make-test-shard-targets-qa-1](../../../../evidence/make-test-shard-targets-qa-1.md)
- verifies [pr-322-ci-takes-11-18-minutes-per-pr](../../../../evidence/pr-322-ci-takes-11-18-minutes-per-pr.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
