# question-c92d1843-1

question · `01M6Q5C2MZF47PAR2KANSXX3WZ`

## session

c92d1843-be8e-4e19-a590-01a1ae548813

## statement

CI-speed change set 01M6Q58 serves no intent yet; the daemon builds only change sets under a milestone. Which intent should it serve? / How should the 271 cmd/rhiza testscripts be split across 4 CI jobs? / After sharding, internal/api (280 s) becomes the critical path, above the ~5 min goal once setup is added. What should CI do with it? / How should the PR keep one verdict for ship-change's pr-green and ci-fix logic?

## Links

### In

- answers [ci-speed-shape](../decision/ci-speed-shape.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
