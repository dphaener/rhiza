# api-snapshot-fixtures-qa-1

evidence · `01M7FJW51F9S9FZWB0S9Y7XDRF` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: ^TestSnapshotFixture passed with -test.count=2 -test.cpu=1,4 shuffled and with -race -test.count=3 -test.cpu=1,4; a mix of serial XDG_CACHE_HOME/run-tmp tests (TestAbandonRemovesRunGoTmp, TestDaemonSweepLeftovers*) with initSnapshot/linkEndpoints users (TestChange*, TestLink*) passed with -test.count=2 -test.cpu=1,4, -test.parallel=1 and -test.parallel=16 shuffled. A -test.timeout=5s kill leaves rhiza-api-cache-*/rhiza-api-snapshots-* behind, and the next run's TestMain sweeps them once over tmpsweep.MaxAge (checked by back-dating them); a normal run leaves none.

## Links

### Out

- verifies [api-snapshot-fixtures](../capability/rhiza/integrity/self-check/api-snapshot-fixtures.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
