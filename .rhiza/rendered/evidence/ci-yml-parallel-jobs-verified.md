# ci-yml-parallel-jobs-verified

evidence · `01M6Q9GC9N5MDP1MPZJNR9A60T` · actor [claude-code-bad97f9a](../actor/claude-code-bad97f9a.md)

## result

PASS

## statement

At 8c0b44ac, go test -trimpath -count=1 -v ./internal/ciworkflow/... passed all 5 tests: the existing triggers, checkout/run-name and concurrency checks plus new TestCIJobsAndMatrices and TestCICheckAggregatesEveryJob. Reading ci.yml and the diff: jobs lint, scripts (matrix shard [1/4..4/4]), api (matrix [1/2, 2/2]), rest, budgets, each with fail-fast false on its matrix, actions/setup-go with cache: true and cache-dependency-path: go.sum, and a run step calling make lint / make test-scripts SHARD / make test-api SHARD / make test-rest / make test-budgets (all exist in the Makefile; make -n test-scripts SHARD=1/4 and test-rest expand as expected); check needs all five, if: always(), and exits 1 unless every needs.*.result is success. PyYAML parses the file to those six jobs. The PR CI duration (about 5 minutes via gh run view) is not observable here: no PR exists for the change set yet, so that point is left to the ship run's CI.

## Links

### Out

- verifies [ci-yml-parallel-jobs](../capability/rhiza/integrity/self-check/ci-yml-parallel-jobs.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
