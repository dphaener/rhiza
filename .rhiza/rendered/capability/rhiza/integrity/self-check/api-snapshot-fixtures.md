# api-snapshot-fixtures

capability · `01M7EQEPVV46QWM31RCD5CVYXW` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./internal/api -run '^TestSnapshotFixture' shows a fixture's copy folds to the same model as building it afresh, a write in one test's copy is not seen by another's, and many goroutines asking for one fixture at once build it once (the build is guarded by sync.Once or the like, so the fixtures are safe for the parallel tests of api-tests-parallel whichever leaf is built first); failure case: a fixture whose build fails fails every test that asks for it with the build's error, not a missing folder. The fixtures are chosen from the profile: go tool pprof -peek 'store.ReadLog' names the helpers and tests whose repeated writes cost most, and those move onto fixtures. Speedup measure: go test -trimpath -count=1 -cpuprofile $TMPDIR/api.prof ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' then go tool pprof -top -cum $TMPDIR/api.prof shows store.ReadLog under 30% of CPU (2026-11-15: 58%). Touches internal/api's test helpers (a new internal/api/fixture_test.go and the helpers in init_test.go, change_test.go and session_test.go) and the _test.go files moved onto the fixtures. Must not change non-test code, must not skip, delete or merge a test, and must not replace FoldLog (projection is out of scope).

## statement

internal/api's tests get their common starting repos (the models many tests rebuild through the same writes) from snapshot fixtures: each is built once per test binary into a template and copied per test, instead of each test re-applying the writes and re-reading the log after each.

## Links

### Out

- decided-by [api-fixture-templates](../../../../decision/api-fixture-templates.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)

### In

- depends-on [api-tests-parallel-finish](api-tests-parallel-finish.md)
- verifies [api-snapshot-fixtures-verify-failed](../../../../evidence/api-snapshot-fixtures-verify-failed.md)
- verifies [api-snapshot-fixtures-verified](../../../../evidence/api-snapshot-fixtures-verified.md)
- verifies [api-snapshot-fixtures-qa-1](../../../../evidence/api-snapshot-fixtures-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
