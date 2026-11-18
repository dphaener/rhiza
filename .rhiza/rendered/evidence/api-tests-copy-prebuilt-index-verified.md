# api-tests-copy-prebuilt-index-verified

evidence · `01M7P7Q2KHRB58PK6Y2NNAPHPS` · actor [claude-code-66862541](../actor/claude-code-66862541.md)

## result

PASS

## statement

Verified at eee179680 under decision api-tests-half-cpu-accept-at-60 (half-CPU line replaced by the measured ~60% cut: builder 193.69 s, earlier verifier 198.84 s user+sys of 321.25 s; profiled run not rerun per the ruling). Checks: block has before/after /usr/bin/time lines and pprof index.Open/catchUp/ensureSchema lines (index.Open 151.85s -> 70.12s). go test -trimpath -count=1 -p 2 -shuffle=11111 ./internal/api: exit 0 (174.1 s); -shuffle=22222: exit 0 (169.4 s). go test -list '.*' ./internal/api counts 839 (>= 838). go test -run TestSnapshotFixture -v: TestSnapshotFixtureCopyCatchesUpFromCopiedIndex PASS (an apply on a copy read back through the copy's index with Read=1, no rebuild; template index size and mtime unchanged, no -wal/-shm). Diff bda3b5d3d..eee179680 touches only internal/api *_test.go and .rhiza: no non-test code, no budget thresholds, TestWritesThroughIndexMatchFold untouched, each copy gets its own index file (templates under a per-binary snapshotDir removed at exit, so no cross-run index), no setup writes baked into the budgets' generated model (recordAda only reuses ctx's own index), no TODOs or skips.

## Links

### Out

- verifies [api-tests-copy-prebuilt-index](../capability/rhiza/integrity/self-check/api-tests-copy-prebuilt-index.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
