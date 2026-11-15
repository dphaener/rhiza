# api-parallel-rest-split-and-finish

decision · `01M7EXM06AB6QZ61P3AB0K9BCN`

## reasoning

The finding's two halves need different remedies, so I took both options it offered. Size: on a recount, the rest set's files outside the session-fixture and runtime groups hold far more tests than the session (214) or runtime (202) leaves, so splitting by prefix brings each half near their size. The session_/web_/inbox_/rex_ files share the daemon-and-pane state likeliest to need serial tests, so they stay together. Scope: a package-wide shuffle, race or timing failure can show up in any group's file, so whichever leaf owns those checks must be allowed to edit every internal/api _test.go file. Putting the checks on a separate last leaf keeps the other leaves' per-file scopes as they are. The intent's limits still hold on every leaf: test code only, no projection, no test skipped, deleted or merged.

## statement

api-tests-parallel-rest is split three ways. Two leaves parallelise its files, each touching only its own files: one takes the session_, web_, inbox_ and rex_ _test.go files that do not call newSessionFixture, plus plan_test.go and lifecycle_test.go (about 140 top-level tests). The other takes the remaining files outside the session and runtime sets: model_, change_, write_, read_, check_, lineage_, guide_, ci_ and the rest (about 195 tests). Each leaf's acceptance is shuffled-twice and -race over its own tests only, with a -run pattern built like the runtime leaf's. A third, finish leaf comes after all four parallelising leaves and may edit any internal/api _test.go file, never non-test code. It owns the package-wide checks: go test -shuffle=on over all of internal/api passes twice with different seeds, -race over the whole package reports no data race, go vet is clean, the test count from go test -list is unchanged, and the 150 s wall-time target holds. It fixes a cross-group race or shuffle failure by making that test serial with a comment, or by fixing its shared state.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ewe7hzgbydtrwyb92zjm65-r3-f2](../question/review-01m7ewe7hzgbydtrwyb92zjm65-r3-f2.md)

### In

- decided-by [api-tests-parallel-model](../capability/rhiza/integrity/self-check/api-tests-parallel-model.md)
- decided-by [api-tests-parallel-finish](../capability/rhiza/integrity/self-check/api-tests-parallel-finish.md)
- decided-by [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
