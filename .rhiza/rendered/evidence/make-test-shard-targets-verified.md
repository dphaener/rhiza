# make-test-shard-targets-verified

evidence · `01M6Q9BBP1J2GJZR42MWE2H6J4` · actor [claude-code-3a599979](../actor/claude-code-3a599979.md)

## result

PASS

## statement

Verifier at 559ffa93: go test -trimpath -count=1 ./cmd/rhiza -run 'TestMake' -v passed TestMakeUpdateDryRun, TestMakeTestTimeout (its code untouched: the diff to makeupdate_test.go only adds lines) and the new TestMakeTestShardTargets, which runs make -n on test-scripts, test-rest, test-budgets and lint and asserts -trimpath plus the budget -skip or -run regex, checks test-api's SHARD=1/2 and 2/2 regexes (via GO_TEST=echo) cover go test -list's internal/api tests exactly once with no budget test, and checks test-budgets prints a fake go's RHIZA_BUDGET_LOG under 'Read-budget timings'. By hand: make test-api SHARD=1/2 GO_TEST=echo prints an anchored ^(...)$ -run with -trimpath; SHARD=3/2 exits non-zero with a usage error; make -n lint runs gofmt -l, go vet and the self-check build; the test target is unchanged. Read the diff: test-rest skips budgets and TestScripts across every package but internal/api, so the four test targets cover make test's tests. Note: in test-api, go test -list's failure is masked by the pipe to sort (no pipefail), so a compile error in internal/api makes a test-api shard report 'no tests' and pass; test-budgets and lint (vet) still fail CI on it.

## Links

### Out

- verifies [make-test-shard-targets](../capability/rhiza/integrity/self-check/make-test-shard-targets.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
