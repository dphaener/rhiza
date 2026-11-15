# api-tests-parallel-rest-qa-1

evidence · `01M7FJW56MGJS29Z67FA3KSBA0` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: The 139 rest-group tests passed shuffled twice (seeds 31337, 424242) and under -race shuffled (106 s) with no data race.

## Links

### Out

- verifies [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
