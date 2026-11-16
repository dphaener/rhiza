# api-tests-on-index

capability · `01M7GFXD744PD29P5BQAMF1287` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api passes, and so does go test -trimpath -count=1 -p 2 -shuffle=on ./internal/api run twice with different seeds; go test -list '.*' ./internal/api counts no fewer tests than before. Profile in the commit's Checks: block: go test -trimpath -count=1 -p 2 -cpuprofile $TMPDIR/api.prof ./internal/api, then go tool pprof -top -cum $TMPDIR/api.prof, before and after, showing store.ReadLog's share of CPU falls (both numbers given). Failure case: a test that counts folds through Ctx.readLog and now sees none is changed to assert on the index instead, not deleted. Touches internal/api/init_test.go and the _test.go files that relied on the fallback. Must not change non-test code (a fix it needs belongs to @parity-interim-from-strict or @parity-write-index-from-strict) and must not share one index between parallel tests.

## statement

internal/api's tests run on the index, as production does: testCtx (internal/api/init_test.go) gives each Ctx an IndexPath under the test's t.TempDir(), and only the tests that exercise the no-index fallback on purpose clear it, each with a comment saying so

## Links

### Out

- depends-on [parity-interim-from-strict](../../views/fast-reads/parity-interim-from-strict.md)
- depends-on [parity-write-index-from-strict](../../views/fast-reads/parity-write-index-from-strict.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- verifies [api-tests-on-index-verified](../../../../evidence/api-tests-on-index-verified.md)
- verifies [api-tests-on-index-qa-1](../../../../evidence/api-tests-on-index-qa-1.md)
- verifies [api-tests-on-index-qa-2](../../../../evidence/api-tests-on-index-qa-2.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
