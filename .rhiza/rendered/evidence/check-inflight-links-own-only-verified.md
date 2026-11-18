# check-inflight-links-own-only-verified

evidence · `01M7PGV9W136ZF1R1APB7AT77P` · actor [claude-code-35e19bf8](../actor/claude-code-35e19bf8.md)

## result

PASS

## statement

At adbc9aeb4: go test -trimpath -count=1 -p 2 ./internal/check/... ok; TestRunInFlightLinks passes and covers an own bad link per set reported once as in-flight, the as-is bad link by the link rule only (once, fail), and the failure case (bad link whose from end the same set retires is still reported), plus equality with the old badLinks-over-every-view rule. Diff touches only internal/check/check.go and check_test.go (plus the actor record); no internal/model change, finding text/rule/severity unchanged. Back to back in the main checkout (load ~5.7): base 845741f94 rhiza check 10.66 s, checker.inFlightLinks 2.16s cum = 14.56% of 14.84s samples; after 8.21 s, inFlightLinks no samples; stdout and stderr findings identical, both exit 0. make self-check in this worktree fails only on uncommitted run state; from a clean detached checkout at adbc9aeb4 it fails only on .rhiza/rendered/index.md lacking actor claude-code-feeecdbd (this leaf's) plus claude-code-376e0ebe and claude-code-a8b56767, stale since base 845741f94, which fails the same way; rhiza render at land fixes it

## Links

### Out

- verifies [check-inflight-links-own-only](../capability/rhiza/integrity/check/check-inflight-links-own-only.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
