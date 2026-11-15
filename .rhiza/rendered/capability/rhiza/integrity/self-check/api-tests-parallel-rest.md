# api-tests-parallel-rest

capability · `01M7EX7WX9045E0A4MD9RZV4A5` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -shuffle=on ./internal/api -run "^($(grep -L newSessionFixture internal/api/session*_test.go internal/api/web*_test.go internal/api/inbox*_test.go internal/api/rex*_test.go internal/api/plan_test.go internal/api/lifecycle_test.go | xargs grep -ho '^func Test[A-Za-z0-9_]*' | cut -c6- | paste -sd'|'))\$" (every test in the session_, web_, inbox_ and rex_ _test.go files that do not call newSessionFixture, plus plan_test.go and lifecycle_test.go: 34 files, 139 top-level tests on 2026-11-15) passes twice in a row with different shuffle seeds; go vet ./internal/api is clean. Failure case: go test -trimpath -race -count=1 ./internal/api with the same -run pattern reports no data race, and a test of these files that swaps a package variable, calls t.Setenv or t.Chdir, sets its own XDG_CACHE_HOME or touches run-tmp (inbox_gate_budget_test.go, web_actions_test.go and the like) stays serial with a comment saying why. Touches only those 34 files. Must not skip, delete or merge a test: go test -list '.*' ./internal/api counts the same tests as before. Must not change non-test code.

## statement

The tests in internal/api's session_, web_, inbox_ and rex_ _test.go files that do not call newSessionFixture, and in plan_test.go and lifecycle_test.go, call t.Parallel() when independent; a test that swaps a package variable (runClock, daemonTicker, daemonReexec and the like), calls t.Setenv or t.Chdir, uses run temp folders, or reads shared machine state stays serial, and a comment on each such test says why it must.

## Links

### Out

- depends-on [api-tests-parallel-session](api-tests-parallel-session.md)
- depends-on [api-tests-parallel-runtime](api-tests-parallel-runtime.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [api-cache-per-binary](../../../../decision/api-cache-per-binary.md)
- decided-by [api-parallel-rest-split-and-finish](../../../../decision/api-parallel-rest-split-and-finish.md)

### In

- depends-on [api-tests-parallel-finish](api-tests-parallel-finish.md)
- blocks [review-01m7ewe7hzgbydtrwyb92zjm65-r3-f2](../../../../question/review-01m7ewe7hzgbydtrwyb92zjm65-r3-f2.md)
- verifies [api-tests-parallel-rest-verified](../../../../evidence/api-tests-parallel-rest-verified.md)
- verifies [api-tests-parallel-rest-qa-1](../../../../evidence/api-tests-parallel-rest-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
