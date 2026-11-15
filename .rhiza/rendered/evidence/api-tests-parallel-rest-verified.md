# api-tests-parallel-rest-verified

evidence · `01M7F733E5ACKFQ744KZKQFS9Z` · actor [claude-code-cc43e589](../actor/claude-code-cc43e589.md)

## result

PASS

## statement

At 4d7d7fb24, the verifier ran the acceptance -run pattern (139 tests from the 34 session_/web_/inbox_/rex_ files without newSessionFixture plus plan_test.go and lifecycle_test.go) with go test -p 2 -trimpath -count=1 under the machine lock, using -shuffle=on, -shuffle=1794751111 and -shuffle=4242: all ok (19-21 s). The same pattern with -race -shuffle=on passed in 100.7 s with no DATA RACE. go vet ./internal/api exited 0. go test -list '.*' ./internal/api lists 756 Test funcs at HEAD and at HEAD~1, so none were skipped, deleted or merged. rhiza ship generate and rhiza render exited 0. The diff only adds lines to the 34 listed _test.go files plus .rhiza actor and rendered records; no non-test code changed. 114 tests call t.Parallel. All 25 serial tests have a // Serial: comment, and each reason holds: limitFixture/useClock/startAPIErr/watchFixture swap package state, launcher.Register is used in rex_actions and rex_install, t.Setenv is used in inbox_gate_budget and webGateFixture, and t.Chdir is used in session_bypass. A scan of the parallel test bodies found no Setenv, Chdir, launcher.Register, runClock, limit or api-error swaps, run-tmp or XDG_CACHE_HOME use.

## Links

### Out

- verifies [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
