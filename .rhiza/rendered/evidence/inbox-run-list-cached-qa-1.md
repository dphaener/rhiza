# inbox-run-list-cached-qa-1

evidence · `01M7PQC1MGDGZ6323Y71MQPG25` · actor [claude-code-f694e280](../actor/claude-code-f694e280.md)

## result

PASS

## statement

QA with $TMPDIR/rhiza-qa built at HEAD 3bf00d3 vs $TMPDIR/rhiza-qabase at base bda3b5d: Scratch clone with XDG_CACHE_HOME in scratch: first HEAD inbox writes run-list/<hash>.gob (13 MB) under the cache dir only; second call 0.39 s. Cache garbage, truncated, empty, a directory in its place, and a read-only cache dir: each run exits 0 with bytes identical to base. Appending a 'stopped: open questions' entry to a journal shows a new open item in HEAD exactly as in base; restoring the journal (shrink) drops it; removing a journal rewrites the cache and matches base; HOME unset/no cache dir runs uncached and matches. go test ./internal/run ok.

## Links

### Out

- verifies [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
