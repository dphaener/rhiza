# review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f1

question · `01M7P39FSX9T4EQV9KQDSKECAK`

## statement

Plan review round 3 left finding f1 (uncovered-intent on inbox-blockers-cheaper): The intent's inbox clause (m39-reads-fast, 200 ms from 645 ms) is not covered by the plan's leaves. The plan's own research profile (m39-test-time-research) puts 37% of a rhiza inbox in run.(*Store).List and 24% in inboxUncoveredProposals, and this leaf names 6.7% for inboxBlockers. That totals 67.7%, so even if all three costs dropped to zero, 32.3% of the base CPU would remain. This leaf's acceptance requires the median at most 31% of the base binary's, and it cannot pass. The cache still stats journals and the uncovered proposals still build their overlays, so the real cuts are smaller than those shares. No leaf serves the rest of the inbox's cost (the remaining ~32%: index open and catch-up, the to-be reader, the other inbox lookups). The critic confirmed it: The research source m39-test-time-research puts 37% of a warm rhiza inbox in run.(*Store).List and 24% in inboxUncoveredProposals, and inbox-blockers-cheaper's statement adds 6.7% for inboxBlockers: 67.7% in all. Even if all three fell to zero, 32.3% would remain, above the 31% (200/645) the acceptance requires, and the real cuts are smaller (the cache still stats journals, uncovered sets still build overlays). No other leaf in the slice (the check-*, whatif-*, api-tests-*, scripts-* and cold-* leaves) targets the rest of inbox's cost. The acceptance cannot pass and the 200 ms clause is not covered: major stands. How should it be settled?

## Links

### Out

- blocks [inbox-blockers-cheaper](../capability/rhiza/orchestration/inbox/inbox-blockers-cheaper.md)

### In

- answers [m39-inbox-rest-leaf](../decision/m39-inbox-rest-leaf.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
