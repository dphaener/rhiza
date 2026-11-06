# ci-yml-parallel-jobs

capability · `01M6Q5HHWA4TNKZ9TB3ZEGH7JR` · parent [self-check](../self-check.md)

## acceptance

A Go test in internal/ciworkflow (go test -trimpath ./internal/ciworkflow/...) parses .github/workflows/ci.yml and asserts the jobs and matrices above, strategy.fail-fast false on each matrix, setup-go with cache on and cache-dependency-path go.sum in every job, each job's run step calling make, and a job named check with needs on every other job, if: always() and a step failing unless every needed result is success; the existing ci-yml triggers and concurrency checks still pass. The change set's own PR CI run completes in about 5 minutes or less, read from gh run view.

## statement

CI runs the suite as parallel jobs of one workflow: lint, a matrix of 4 script shards, a matrix of 2 internal/api shards, rest and budgets, each calling its make target with fail-fast off and actions/setup-go caching modules and the build cache keyed on go.sum, and an aggregate job named check that needs every other job, runs if: always() and fails unless each succeeded; a PR's CI gives its verdict in about 5 minutes.

## Links

### Out

- depends-on [make-test-shard-targets](make-test-shard-targets.md)
- decided-by [ci-speed-shape](../../../../decision/ci-speed-shape.md)

### In

- verifies [ci-yml-parallel-jobs-verified](../../../../evidence/ci-yml-parallel-jobs-verified.md)
- verifies [ci-yml-parallel-jobs-qa-1](../../../../evidence/ci-yml-parallel-jobs-qa-1.md)
- verifies [ci-yml-parallel-jobs-fix-2-verified](../../../../evidence/ci-yml-parallel-jobs-fix-2-verified.md)
- verifies [pr-322-ci-takes-11-18-minutes-per-pr](../../../../evidence/pr-322-ci-takes-11-18-minutes-per-pr.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
