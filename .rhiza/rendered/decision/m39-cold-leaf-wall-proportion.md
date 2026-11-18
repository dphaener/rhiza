# m39-cold-leaf-wall-proportion

decision · `01M7P3DVMNDXBE0BQRDTMW7ZRF`

## reasoning

The finding is right: m39-cold is wall time, and the base log shows wall does not follow CPU (400 s wall for 963 CPU-s on 4 cores, cmd/rhiza alone 254 s wall), so a CPU ratio can pass while the suite misses 240 s. The fix keeps the plan's proportion style (the inbox and check leaves compare commit to base back to back, since the shared host is noisy and an absolute number would vary with load) but applies it to what the target measures: the wall of the three slowest packages and cmd/rhiza's wall, the critical path the leaf's statement claims to shorten. 60% matches 240/400 and lowers no target. It does not require make test in the run, since AGENTS.md forbids agents in runs from running make test; the real cold make test is measured as m39-cold says once the change set has landed.

## statement

cold-run-ship-scripts-cut's combined acceptance measures wall time, not only CPU: in the back-to-back runs of /usr/bin/time -f '%e s %U user %S sys' go test -trimpath -count=1 ./internal/api ./cmd/rhiza ./internal/index at the commit and in a git worktree at the base 5d2596d, the commit's wall time is at most 60% (240 s / 400 s) of the base's, and cmd/rhiza's own package time (its go test ok line) is at most 60% of the base's; the user+sys CPU ratio stays as a further check, and both runs' wall, CPU and per-package times go in the Checks: block. The full cold make test against m39-cold's 240 s stays the milestone's own measurement, taken by the orchestrator after the change set lands, not by the run

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f3](../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f3.md)

### In

- decided-by [cold-run-ship-scripts-cut](../capability/rhiza/integrity/self-check/cold-run-ship-scripts-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
