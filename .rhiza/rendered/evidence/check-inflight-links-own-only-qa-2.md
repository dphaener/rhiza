# check-inflight-links-own-only-qa-2

evidence · `01M7PRNFFZJVTZRGFAA5NDBDN0` · actor [claude-code-6fe9e43d](../actor/claude-code-6fe9e43d.md)

## result

PASS

## statement

Built $TMPDIR/rhiza-qa at HEAD 639e806f6 and $TMPDIR/rhiza-qa-base at merge-base bda3b5d3d (go build -trimpath, never installed). rhiza check --json in the main checkout: both exit 0, identical bytes (user+sys 7.2 s HEAD vs 16.4 s base). In a scratch clone with the main checkout's .rhiza and runs, after rhiza note add and rhiza accept --rule owner entries, both binaries' check --json again identical, exit 0

## Links

### Out

- verifies [check-inflight-links-own-only](../capability/rhiza/integrity/check/check-inflight-links-own-only.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
