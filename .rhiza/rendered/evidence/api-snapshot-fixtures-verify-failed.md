# api-snapshot-fixtures-verify-failed

evidence · `01M7FBEGAE5HH3TXGVGSQ0GDCX` · actor [claude-code-790f01c5](../actor/claude-code-790f01c5.md)

## result

FAIL

## statement

Verifier at 8ade328e8: FAIL. Defect: internal/api/fixture_test.go snapshotBase/make builds the template under os.Getenv("XDG_CACHE_HOME") read at build time, but serial tests t.Setenv XDG_CACHE_HOME to their own t.TempDir() and then call initCtx (sliceCtx in slice_test.go, TestGateAnswerUsesWorktreeIndex and inbox_gate_budget_test.go:187, web_actions_test.go:66). If one triggers the first initSnapshot build, the template is deleted with that test's TempDir and every later copy fails. Reproduced: go test -p 2 -trimpath -count=1 ./internal/api -run '^(TestSliceViews|TestSnapshotFixtureCopiesAreIsolated|TestChangeOutsideRepo)$' -> FAIL 'snapshot init: copy: lstat /home/rhiza/tmp/TestSliceViews.../001/snapshot-...: no such file or directory'. Order-dependent (-shuffle, -run subsets). Fix: choose the template base once per binary in TestMain (a package var beside the rhiza-api-cache- folder, removed at exit), not from the env at build time, plus a regression test where a build under a t.Setenv(XDG_CACHE_HOME) test still serves later copies. Passing: Checks 1 and 2 (race+shuffle) exit 0; go vet exit 0; -list 772 entries (Test funcs 756 -> 760); profile run exit 1 only TestDaemonProfile/TestDaemonProfileDefaultPath, same at 8ade328e8^ (pre-existing); store.ReadLog 25.07% cum, package 151.8 s; TestLinkEndpointKinds keeps every case and assertion; no non-test code changed. rhiza run findings --leaf was refused (no leaves step in a conductor run), so this is recorded as Evidence.

## Links

### Out

- verifies [api-snapshot-fixtures](../capability/rhiza/integrity/self-check/api-snapshot-fixtures.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
