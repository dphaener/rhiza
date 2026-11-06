# ci-speed-plan-interview

source · `01M6Q5GZSKS19EP83EA3PDTN5S`

## quotes

```json
[]
```

## raw_sha256

c35bf3b3aea59477fa04b560ecbf4339d602aefcb2024a51f6304145df9090d1

## session

c92d1843-be8e-4e19-a590-01a1ae548813

## summary

Autonomous plan run for change set 01M6Q58 (CI takes 11-18 minutes per PR). The planner read ci.yml, the Makefile's test target, cmd/rhiza's TestScripts, the budget tests and ship check's rollup parsing, and raised four choices as one question: the intent to serve, how to shard the testscripts, what to do with internal/api once it becomes the critical path, and how the PR keeps one verdict. The arbiter answered with Decision ci-speed-shape: serve milestone-25; shard TestScripts round-robin by RHIZA_TEST_SHARD=i/n over the sorted script names; shard internal/api's top-level tests into 2 jobs by an anchored -run regex; a rest job, a budgets job and a lint job; Makefile targets per job so CI calls make; one workflow with fail-fast false, setup-go caching keyed on go.sum, and an aggregate job named check under if: always() that fails unless every job succeeded; ship check unchanged unless a test shows otherwise. No owner words were typed in this run.

## title

CI speed plan, autonomous interview with the arbiter

## Links

### In

- derived-from [ci-speed-shape](../decision/ci-speed-shape.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
