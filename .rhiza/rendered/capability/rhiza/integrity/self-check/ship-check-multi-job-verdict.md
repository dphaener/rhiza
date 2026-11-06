# ship-check-multi-job-verdict

capability · `01M6Q5HHWYFNJV9D7A7KNYSZGH` · parent [self-check](../self-check.md)

## acceptance

Unit tests in internal/api (go test -trimpath ./internal/api -run 'TestParsePRView|TestChecksVerdict') feed a statusCheckRollup of the lint, shard, rest, budgets and check jobs of one run with one failed script shard and a cancelled or failed aggregate, and assert the verdict is settled and not green, failed names that shard and the aggregate, and every check carries the same RunID; a rollup where all jobs succeeded is green; one still in progress is not settled.

## statement

ship check gives one verdict over a PR whose CI is several jobs of one run: green only when every job and the aggregate check completed and passed, and a failed shard is named with the run ID its log is read from.

## Links

### Out

- decided-by [ci-speed-shape](../../../../decision/ci-speed-shape.md)

### In

- verifies [ship-check-multi-job-verdict-verified](../../../../evidence/ship-check-multi-job-verdict-verified.md)
- verifies [ship-check-multi-job-verdict-qa-1](../../../../evidence/ship-check-multi-job-verdict-qa-1.md)
- verifies [pr-322-ci-takes-11-18-minutes-per-pr](../../../../evidence/pr-322-ci-takes-11-18-minutes-per-pr.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
