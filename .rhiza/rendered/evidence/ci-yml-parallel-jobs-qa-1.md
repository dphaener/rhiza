# ci-yml-parallel-jobs-qa-1

evidence · `01M6QB1P7M9GG7A21T35EV5EFJ` · actor [claude-code-49d75309](../actor/claude-code-49d75309.md)

## result

PASS

## statement

QA at 2f747105: go test -trimpath -count=1 ./internal/ciworkflow/... passed. Read ci.yml: lint, scripts [1/4..4/4], api [1/2,2/2], rest, budgets, fail-fast false, setup-go cache on go.sum, and a check job (needs all five, if: always(), exits 1 on any non-success). Each make target it calls exists and ran or expanded locally (make -n lint, make -n test-rest). The about-5-minutes PR CI point is not observable yet: no PR or Actions run exists for the change set, so it falls to the ship run's CI.

## Links

### Out

- verifies [ci-yml-parallel-jobs](../capability/rhiza/integrity/self-check/ci-yml-parallel-jobs.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
