# ship-check-multi-job-verdict-qa-1

evidence · `01M6QB1P5Q902WD3NT5EXY907R` · actor [claude-code-49d75309](../actor/claude-code-49d75309.md)

## result

PASS

## statement

QA at 2f747105: go test -trimpath -count=1 ./internal/api -run 'TestParsePRView|TestChecksVerdict' passed; ship.go is unchanged in the diff, so the new tests pin the existing verdict over a multi-job rollup. Not exercised against a real PR: none exists yet.

## Links

### Out

- verifies [ship-check-multi-job-verdict](../capability/rhiza/integrity/self-check/ship-check-multi-job-verdict.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
