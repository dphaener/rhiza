# review-01m7ewe7hzgbydtrwyb92zjm65-r3-f2

question · `01M7EXHX2H6FVD0YE2H4KJB5HP`

## statement

Plan review round 3 left finding f2 (oversized-leaf on api-tests-parallel-rest): The leaf does two jobs. It parallelises 86 _test.go files: 335 top-level tests, ~18.8k lines, against 202 for runtime and 214 for session. It also owns the package-wide finish for all three leaves: the whole package passes shuffled twice, -race passes over the whole package, and the 150 s wall-time target holds. Yet it 'touches only internal/api _test.go files outside the other two leaves' sets', so a shuffle or race failure that shows up in session or runtime files when everything runs together cannot be fixed inside this leaf. That is too much for one fresh session. Split the rest files (e.g. web_/inbox_/rex_/session_* without the fixture vs model/change/write/read and the rest), or let the package-wide shuffle/race/timing check fix tests in any file. The critic confirmed it: Verified counts: the rest set is 86 _test.go files, 335 top-level tests, 18,838 lines, against session's 24 files/214 tests/9,003 lines and runtime's 65 files/202 tests/13,448 lines. On top of that the rest leaf owns the package-wide finish (shuffled twice over all of internal/api, -race over the whole package, the 150 s wall target) while its scope says it touches only files outside the other two leaves' sets, so a cross-group race or shuffle failure surfacing in session or runtime files cannot be fixed within the leaf. Too much for one session and partly unbuildable as scoped; oversized-leaf major stands. How should it be settled?

## Links

### Out

- blocks [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)

### In

- answers [api-parallel-rest-split-and-finish](../decision/api-parallel-rest-split-and-finish.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
