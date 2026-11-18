# inbox-blockers-cheaper-qa-2

evidence · `01M7PRNGBD5RXADGYZNP23N1H2` · actor [claude-code-6fe9e43d](../actor/claude-code-6fe9e43d.md)

## result

PASS

## statement

Built $TMPDIR/rhiza-qa at HEAD 639e806f6 and $TMPDIR/rhiza-qa-base at merge-base bda3b5d3d. rhiza inbox --json warm in the main checkout: identical bytes, 10 items (HEAD 0.56 s vs base 1.0 s warm). In a scratch clone of the checkout with its .rhiza and run journals, under clean env -i: after rhiza note add the items are unchanged and identical on both; after rhiza accept --rule owner --change 01M7NVMG0EJFRJV03G025C6366 both list 9 items, identical bytes, so no stale index or view is served

## Links

### Out

- verifies [inbox-blockers-cheaper](../capability/rhiza/orchestration/inbox/inbox-blockers-cheaper.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
