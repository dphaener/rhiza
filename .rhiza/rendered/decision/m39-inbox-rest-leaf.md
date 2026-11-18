# m39-inbox-rest-leaf

decision · `01M7P3AWHGFS1PZTAZ84R3TH5W`

## reasoning

The finding's arithmetic holds: m39-test-time-research puts 37% of a warm inbox in run.(*Store).List and 24% in inboxUncoveredProposals, and inbox-blockers-cheaper names 6.7% for inboxBlockers, so even if all three fell to zero 32.3% would remain, above the 31% its acceptance asks; no leaf can pass a bar its own cuts cannot reach. The intent says never lower a target, so relaxing the proportion is out, and it says profile before changing anything, which m39-time-profile-first applies per leaf. The remaining cost is not yet profiled in detail (the research names only the top three), so the leaf that serves it profiles after the three cuts have landed, when what is left is what it measures, and the combined 200/645 check belongs on that last leaf, where every inbox cut is in. This covers the uncovered intent clause with the least change to the existing leaves.

## statement

The inbox clause of m39-reads-fast (200 ms, from 645 ms) gets a fourth inbox leaf that cuts the rest of a warm rhiza inbox's cost: the ~32% the three leaves leave (index open and catch-up, the to-be reader, the other inbox lookups), starting from a fresh RHIZA_CPUPROFILE of rhiza inbox taken once inbox-run-list-cached, inbox-proposals-overlay-once and inbox-blockers-cheaper have landed, and cutting the top costs it names in internal/api, internal/index or internal/model. That leaf depends on the three and carries the combined acceptance (median user+sys CPU of a warm rhiza inbox at most 200/645 of the base binary's at 5d2596d, measured as inbox-blockers-cheaper now states it). inbox-blockers-cheaper keeps only its own share: its acceptance checks inboxBlockers', planGather's and planCore's pprof -top -cum share falling and its failure case, and drops the 31% proportion. The 200 ms target is not lowered.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f1](../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f1.md)

### In

- decided-by [inbox-rest-cut](../capability/rhiza/orchestration/inbox/inbox-rest-cut.md)
- decided-by [inbox-blockers-cheaper](../capability/rhiza/orchestration/inbox/inbox-blockers-cheaper.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
