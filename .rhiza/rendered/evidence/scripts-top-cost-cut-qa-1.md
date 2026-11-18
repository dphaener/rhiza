# scripts-top-cost-cut-qa-1

evidence · `01M7PQC1QNJ8AYT4PFK35JQ66P` · actor [claude-code-f694e280](../actor/claude-code-f694e280.md)

## result

PASS

## statement

QA with $TMPDIR/rhiza-qa built at HEAD 3bf00d3 vs $TMPDIR/rhiza-qabase at base bda3b5d: go test ./cmd/rhiza -run TestScripts/run-ship ok under the changed harness (GC/one-P settings in cmd/rhiza/main_test.go), RHIZA_/SHIP_/REX_ unset; timing pair read from the commit's Checks: block, not rerun.

## Links

### Out

- verifies [scripts-top-cost-cut](../capability/rhiza/integrity/self-check/scripts-top-cost-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
