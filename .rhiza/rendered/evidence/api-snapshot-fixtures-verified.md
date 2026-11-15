# api-snapshot-fixtures-verified

evidence · `01M7FC5DX791GWBA6KBXWGYAKX` · actor [claude-code-ba79a679](../actor/claude-code-ba79a679.md)

## result

PASS

## statement

Verified at d41a8ddb7 (8ade328e8 plus the fix d41a8ddb7). The fix diff is test code only (budget_test.go TestMain makes a rhiza-api-snapshots- folder into snapshotDir, sweeps stale ones and removes it at exit; fixture_test.go snapshotBase returns snapshotDir and no longer reads XDG_CACHE_HOME). Builds read no other fixture state from the environment. The leaf's non-test changes are only .rhiza/log and .rhiza/rendered records. The previous finding's reproducer go test -p 2 -trimpath -count=1 ./internal/api -run '^(TestSliceViews|TestSnapshotFixtureCopiesAreIsolated|TestChangeOutsideRepo)$' exits 0. go test -race -shuffle=on -run '^TestSnapshotFixture' -v passed all 5 tests, including the new TestSnapshotFixtureBuildUnderSetenvServesLaterCopies. The Checks: shuffle run over TestSlice|TestInboxGateQuestions|TestGateAnswerUsesWorktreeIndex|TestWeb|TestSnapshotFixture|TestChangeOutsideRepo exits 0. go vet exits 0 and gofmt -l prints nothing. go test -list counts 761, as the Checks: block says. The brief's 768 does not match this count, but the top-level test functions differ from 8ade328e8~1 only by the 5 added TestSnapshotFixture* tests: none skipped, deleted or merged, and no new t.Skip or TODO. Profile measure (-cpuprofile with the acceptance's -skip, -o $TMPDIR/api.test): store.ReadLog is 24.37% cum, under 30%. In that run TestDaemonProfile and TestDaemonProfileDefaultPath fail with 'cpu profiling already in use': they start a real CPU profile, which clashes with -cpuprofile. Neither file was touched by this leaf, and both tests pass without -cpuprofile (exit 0).

## Links

### Out

- verifies [api-snapshot-fixtures](../capability/rhiza/integrity/self-check/api-snapshot-fixtures.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
