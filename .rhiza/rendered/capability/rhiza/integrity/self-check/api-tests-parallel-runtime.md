# api-tests-parallel-runtime

capability · `01M7EX7WWQZT7JZK1MD30V9086` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -shuffle=on ./internal/api -run "^($(grep -L newSessionFixture internal/api/daemon_*_test.go internal/api/run_*_test.go internal/api/ship_*_test.go internal/api/watchdog_*_test.go internal/api/health_*_test.go internal/api/usage_*_test.go internal/api/feed_*_test.go | xargs grep -ho '^func Test[A-Za-z0-9_]*' | cut -c6- | paste -sd'|'))\$" (every test in those files) passes twice in a row with different shuffle seeds; go vet ./internal/api is clean. Failure case: go test -trimpath -race -count=1 ./internal/api with the same -run pattern (the daemonTicker, daemonReexec, webFactory, pollGitTimeout and runClock swappers, run_handback_test.go's and run_spawned_test.go's tests, beside the parallel ones) reports no data race. Touches only the _test.go files of those prefixes that do not call newSessionFixture. Must not skip, delete or merge a test: go test -list '.*' ./internal/api counts the same tests as before. Must not change non-test code.

## statement

The tests in internal/api's daemon_, run_, ship_, watchdog_, health_, usage_ and feed_ _test.go files that do not call newSessionFixture call t.Parallel() when independent; a test that swaps a package variable (runClock, daemonTicker, daemonReexec and the like), calls t.Setenv or t.Chdir, uses run temp folders, or reads shared machine state stays serial, and a comment on each such test says why it must.

## Links

### Out

- depends-on [api-tests-parallel-session](api-tests-parallel-session.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [api-cache-per-binary](../../../../decision/api-cache-per-binary.md)

### In

- depends-on [api-tests-parallel-rest](api-tests-parallel-rest.md)
- depends-on [api-tests-parallel-model](api-tests-parallel-model.md)
- depends-on [api-tests-parallel-finish](api-tests-parallel-finish.md)
- verifies [api-tests-parallel-runtime-verified](../../../../evidence/api-tests-parallel-runtime-verified.md)
- verifies [api-tests-parallel-runtime-qa-1](../../../../evidence/api-tests-parallel-runtime-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
