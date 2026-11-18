# review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f3

question · `01M7P39FYE2F2ZH7JD8XJVHSRV`

## statement

Plan review round 3 left finding f3 (missing-acceptance on cold-run-ship-scripts-cut): m39-cold is the wall time of a cold make test (240 s), but this acceptance checks a CPU proportion instead: user+sys of internal/api, cmd/rhiza and internal/index at most 240/400 of the base's. In the base cold log, wall time does not follow CPU: 400 s wall for 963 CPU-s on 4 cores, and cmd/rhiza alone took 254 s wall, more than the whole 240 s target. The leaf's statement says the suite's critical path shortens, yet nothing measures cmd/rhiza's own wall time or make test's wall time. The CPU ratio can pass while the cold suite still misses 240 s. The critic confirmed it: The research source gives the cold make test as 400 s wall with cmd/rhiza alone at 254 s wall and notes that package times overlap under contention, so wall time does not follow CPU. cold-run-ship-scripts-cut's acceptance compares only user+sys CPU of the three packages to 240/400 of the base's and merely reports wall times; nothing checks cmd/rhiza's wall time or make test's wall against m39-cold's 240 s. The CPU ratio can pass while the cold suite still misses its wall-time target, so the acceptance does not measure the target the statement claims: major stands. How should it be settled?

## Links

### Out

- blocks [cold-run-ship-scripts-cut](../capability/rhiza/integrity/self-check/cold-run-ship-scripts-cut.md)

### In

- answers [m39-cold-leaf-wall-proportion](../decision/m39-cold-leaf-wall-proportion.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
