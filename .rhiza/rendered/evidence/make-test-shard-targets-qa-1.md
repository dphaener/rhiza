# make-test-shard-targets-qa-1

evidence · `01M6QB237C6PKPWK5RH1YKF014` · actor [claude-code-49d75309](../actor/claude-code-49d75309.md)

## result

PASS

## statement

QA at 2f747105, real runs on the VM: make test-scripts SHARD=2/4 passed (68 scripts); make test-api SHARD=2/2 passed (ok internal/api 267 s on a loaded machine); make test-budgets passed for internal/api, index and slice and printed the Read-budget timings block; make test-api SHARD=0/4, 5/4, x and 1/0 each exit 1 with 'SHARD=...: want i/n with 1 <= i <= n'. With GO_TEST=echo, shards 1/2 and 2/2 list 229 and 228 tests, disjoint, covering all 460 top-level internal/api tests except the 3 budget tests. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScriptShard|TestMake' passed; the test target is unchanged in the diff.

## Links

### Out

- verifies [make-test-shard-targets](../capability/rhiza/integrity/self-check/make-test-shard-targets.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
