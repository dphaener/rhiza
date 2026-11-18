# check-log-decode-cheaper-verified

evidence · `01M7PH944K7VY618ZKYMQYJJY2` · actor [claude-code-b384e89d](../actor/claude-code-b384e89d.md)

## result

PASS

## statement

Verified at f865e9a2c (verifier, attempt 1). go test -trimpath -count=1 -p 2 ./internal/store/... ./internal/entry/...: exit 0. Failure case: TestReadEntriesBadFilesSameErrors (non-canonical, unparsable and canonical-but-invalid entry files) passes with the parent adbc9aeb4's entry.go too, so errors are unchanged. make self-check in a detached clean worktree at f865e9a2c: exit 0. RHIZA_CPUPROFILE of rhiza check --json in the main checkout, back to back at load ~5: before (adbc9aeb4) 7.96 s wall, 9.44 user + 2.01 sys = 11.45 s CPU; pprof -top -cum store.readEntries.func1 15.45%, store.readEntries 3.88%, entry.Decode 12.10%, canon.Canonicalize 11.02% of 11.07 s. After 7.31 s wall, 7.77 user + 1.93 sys = 9.70 s CPU; store.readEntries.func1 9.61%, store.readEntries 2.38%, entry.Decode 5.29%, canon.Canonicalize 8.96% of 9.26 s. readEntries share fell. Before and after rhiza check --json: same bytes (cmp), both exit 0. Diff touches only internal/entry/entry.go and internal/store/read_test.go plus .rhiza records; Decode still runs canon.Check first and every envelope rule; existing tests unchanged; no TODOs.

## Links

### Out

- verifies [check-log-decode-cheaper](../capability/rhiza/integrity/check/check-log-decode-cheaper.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
