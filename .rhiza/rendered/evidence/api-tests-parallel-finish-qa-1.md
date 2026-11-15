# api-tests-parallel-finish-qa-1

evidence · `01M7FJW5B2BTFRJVZSZEWZYN97` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: The whole package, skipping the three budget tests, passed with -test.shuffle=8086 in 157 s wall on the busy shared VM and left no new temp folder; go vet ./internal/api is clean; the top-level Test funcs differ from the merge-base only by the five added TestSnapshotFixture* tests and TestMain (none removed); non-test .go files are unchanged from the merge-base.

## Links

### Out

- verifies [api-tests-parallel-finish](../capability/rhiza/integrity/self-check/api-tests-parallel-finish.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
