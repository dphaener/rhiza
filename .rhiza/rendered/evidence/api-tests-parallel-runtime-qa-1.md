# api-tests-parallel-runtime-qa-1

evidence · `01M7FJW54XSSZXGZ4EPHBJK7AX` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: The 208 runtime-file tests passed shuffled twice (seeds 31337, 424242; ~65 s each); under -race the whole group overran a 250 s timeout as one command, so it ran in three parts (daemon_ 46, run_ 79, ship_/watchdog_/health_/usage_/feed_ 83 tests; 53 s, 28 s, 210 s), each passing with no data race. The slow part is the serial usage_ tests at 23-28 s each under -race, not a hang.

## Links

### Out

- verifies [api-tests-parallel-runtime](../capability/rhiza/integrity/self-check/api-tests-parallel-runtime.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
