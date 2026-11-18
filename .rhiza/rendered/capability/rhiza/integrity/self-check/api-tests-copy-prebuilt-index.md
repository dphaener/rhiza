# api-tests-copy-prebuilt-index

capability · `01M7P20X412DYC423YKYB9Y917` · parent [self-check](../self-check.md)

## acceptance

Checks: block gives, before and after, /usr/bin/time -f '%e s %U user %S sys' go test -trimpath -count=1 -p 2 -cpuprofile $TMPDIR/api.prof ./internal/api and the index-setup lines of go tool pprof -top -cum $TMPDIR/api.prof (index.Open, ensureSchema, catch-up), the after run's user+sys CPU at most half the before's; go test -trimpath -count=1 -p 2 -shuffle=on ./internal/api passes twice with different seeds; go test -list '.*' ./internal/api counts no fewer tests than before. Failure case: a new test in internal/api/fixture_test.go writes an entry (an apply) on a snapshot copy and reads it back through the copy's index, which catches up from the copied index file; and the template's own index file is unchanged after that test (its size and modification time compared), so no copy writes the template's index

## files

internal/api/fixture_test.go, internal/api/init_test.go, internal/api/*_test.go that build snapshot templates or call testCtx on a snapshot copy

## must_not

Change non-test code, including internal/index's freshness check or SQLite pragmas; share one index file or open index between two tests; lower a budget test threshold or delete, skip or shorten a test; regenerate internal/api/testdata/test-timings.tsv by hand

## statement

internal/api's tests no longer catch each test's index up from an empty file: every snapshot template carries the index its build caught up, snapshot.copy gives each copy its own copy of that index file under the test's t.TempDir(), and the index setup left per test in the profile (open, schema, catch-up) is cut

## Links

### Out

- decided-by [m39-api-index-per-test-copy](../../../../decision/m39-api-index-per-test-copy.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [cold-run-ship-scripts-cut](cold-run-ship-scripts-cut.md)
- blocks [api-tests-half-cpu-unreachable](../../../../question/api-tests-half-cpu-unreachable.md)
- verifies [api-tests-copy-prebuilt-index-verified](../../../../evidence/api-tests-copy-prebuilt-index-verified.md)
- verifies [api-tests-copy-prebuilt-index-qa-1](../../../../evidence/api-tests-copy-prebuilt-index-qa-1.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
