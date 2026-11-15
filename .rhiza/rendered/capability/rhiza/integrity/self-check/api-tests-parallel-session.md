# api-tests-parallel-session

capability · `01M7EX7WW01FJ4NQPJ2GPYEZNX` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -shuffle=on ./internal/api -run "^($(grep -l newSessionFixture internal/api/*_test.go | xargs grep -ho '^func Test[A-Za-z0-9_]*' | cut -c6- | paste -sd'|'))\$" (every test in the files that call newSessionFixture) passes twice in a row with different shuffle seeds, and the same run with -race reports no data race; grep shows no t.Setenv left in newSessionFixture and TestMain in internal/api/budget_test.go sets XDG_CACHE_HOME and removes the folder when the binary exits; go vet ./internal/api is clean. Failure case: a session-fixture test that starts a run step or inspects run-tmp stays serial with a comment saying why, and no test writes under the host's ~/.cache. Touches internal/api/budget_test.go, session_test.go's newSessionFixture and the _test.go files that call newSessionFixture only. Must not skip, delete or merge a test: go test -list '.*' ./internal/api counts the same tests as before. Must not change non-test code.

## statement

internal/api's TestMain (in budget_test.go) points XDG_CACHE_HOME at one temporary folder for the whole test binary and newSessionFixture stops calling t.Setenv, and the tests in the _test.go files that call newSessionFixture call t.Parallel(); among them a test that swaps a package variable (runClock, daemonTicker, daemonReexec and the like), calls t.Setenv or t.Chdir, uses run temp folders, or reads shared machine state stays serial, and a comment on each such test says why it must.

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [api-cache-per-binary](../../../../decision/api-cache-per-binary.md)

### In

- depends-on [api-tests-parallel-runtime](api-tests-parallel-runtime.md)
- depends-on [api-tests-parallel-rest](api-tests-parallel-rest.md)
- depends-on [api-tests-parallel-model](api-tests-parallel-model.md)
- depends-on [api-tests-parallel-finish](api-tests-parallel-finish.md)
- verifies [api-tests-parallel-session-verified](../../../../evidence/api-tests-parallel-session-verified.md)
- verifies [api-tests-parallel-session-qa-1](../../../../evidence/api-tests-parallel-session-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
