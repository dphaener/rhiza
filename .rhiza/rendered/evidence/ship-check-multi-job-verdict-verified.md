# ship-check-multi-job-verdict-verified

evidence · `01M6Q8Z192G8ENXPBKQ46CH0HC` · actor [claude-code-43f7013f](../actor/claude-code-43f7013f.md)

## result

PASS

## statement

Verifier claude-code-43f7013f at 04879199: ran go test -trimpath -count=1 -v ./internal/api -run 'TestParsePRView|TestChecksVerdict'; all PASS incl. new TestParsePRViewMultiJob. Read the diff (internal/api/ship_test.go plus the builder's actor record only): it feeds parsePRView a statusCheckRollup of lint, scripts (1/4) FAILURE, scripts (2/4), rest, budgets and check CANCELLED, all with detailsUrl under actions/runs/4242, and asserts settled, not green, failed == [scripts (1/4), check], and every check's RunID == 4242; an all-SUCCESS rollup is settled and green with no failures; a rollup with IN_PROGRESS/QUEUED jobs is neither settled nor green. Matches checksVerdict and actionsRunID in internal/api/ship.go.

## Links

### Out

- verifies [ship-check-multi-job-verdict](../capability/rhiza/integrity/self-check/ship-check-multi-job-verdict.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
