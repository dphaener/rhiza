# api-tests-parallel-model

capability · `01M7F32H400Z0D8A1WY3GKTD3S` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -shuffle=on ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' -run "^($(ls internal/api/*_test.go | grep -v -e /session -e /web -e /inbox -e /rex -e /plan_test -e /lifecycle_test -e /daemon_ -e /run_ -e /ship_ -e /watchdog_ -e /health_ -e /usage_ -e /feed_ | xargs grep -L newSessionFixture | xargs grep -ho '^func Test[A-Za-z0-9_]*' | cut -c6- | paste -sd'|'))\$" (every test in the internal/api _test.go files outside the session, runtime and api-tests-parallel-rest sets: model_, change_, write_, read_, check_, lineage_, guide_, ci_ and the rest, 52 files and 196 top-level tests on 2026-11-15) passes twice in a row with different shuffle seeds; go vet ./internal/api is clean. Failure case: go test -trimpath -race -count=1 ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' with the same -run pattern reports no data race, and a test of these files that swaps a package variable, calls t.Setenv or t.Chdir, sets its own XDG_CACHE_HOME or touches run-tmp (slice_test.go and the like) stays serial with a comment saying why. Touches only those files (budget_test.go's TestMain belongs to api-tests-parallel-session). Must not skip, delete or merge a test: go test -list '.*' ./internal/api counts the same tests as before. Must not change non-test code.

## statement

The tests in the internal/api _test.go files that neither api-tests-parallel-session, api-tests-parallel-runtime nor api-tests-parallel-rest covers (model_, change_, write_, read_, check_, lineage_, guide_, ci_ and the rest) call t.Parallel() when independent; a test that swaps a package variable, calls t.Setenv or t.Chdir, uses run temp folders, or reads shared machine state stays serial, and a comment on each such test says why it must.

## Links

### Out

- depends-on [api-tests-parallel-session](api-tests-parallel-session.md)
- depends-on [api-tests-parallel-runtime](api-tests-parallel-runtime.md)
- decided-by [api-parallel-rest-split-and-finish](../../../../decision/api-parallel-rest-split-and-finish.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [api-cache-per-binary](../../../../decision/api-cache-per-binary.md)

### In

- depends-on [api-tests-parallel-finish](api-tests-parallel-finish.md)
- verifies [api-tests-parallel-model-verified](../../../../evidence/api-tests-parallel-model-verified.md)
- verifies [api-tests-parallel-model-qa-1](../../../../evidence/api-tests-parallel-model-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
