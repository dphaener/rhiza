# inbox-rest-cut-qa-1

evidence · `01M7PQC1H9Y5S6Z1AFDMRQBAT8` · actor [claude-code-f694e280](../actor/claude-code-f694e280.md)

## result

PASS

## statement

QA with $TMPDIR/rhiza-qa built at HEAD 3bf00d3 vs $TMPDIR/rhiza-qabase at base bda3b5d: Warm rhiza inbox in the main checkout: base 1.36/1.22 s (1.35/1.32 u), HEAD 0.35/0.46 s (0.31 u), same bytes. Scratch clone: after rhiza note add and after rhiza accept, HEAD and base list the same items, so no stale index or view is served.

## Links

### Out

- verifies [inbox-rest-cut](../capability/rhiza/orchestration/inbox/inbox-rest-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
