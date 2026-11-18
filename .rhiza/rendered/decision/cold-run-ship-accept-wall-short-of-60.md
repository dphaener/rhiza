# cold-run-ship-accept-wall-short-of-60

decision · `01M7PGB8YHWT5EJZVNGVBY71V9`

## reasoning

This follows the slice's two precedents, api-tests-half-cpu-accept-at-60 and scripts-cpu-cut-accept-short-of-15: a per-leaf proportion line is a proxy, the intent's 'never lower a target' binds the milestone target, budget thresholds and tests, and accepting this leaf at a measured miss lowers none of them. m39-cold-leaf-wall-proportion made the 60% check wall time so the miss would show before landing, and it does show: the builder and verifier agree within 1% on every figure. The question is right that no cut inside the leaf's files and must_not can reach 60%: internal/api alone is 370 s of the 372 s wall and the leaf may not edit it; cmd/rhiza's CPU is in about 107k child processes, the harness levers (GOMAXPROCS=1, git housekeeping off, GOGC/GOMEMLIMIT) are spent, strace shows no I/O waits left, and even perfect packing of its CPU on 4 cores misses 138.6 s. What is left is production code that scripts-cpu-cut-accept-short-of-15 already ruled out of a test-harness leaf. A third attempt could only fail again or break must_not. On the commit: the acceptance's single before/after pair on a host at load 5-9 cannot resolve a few-percent cut either way, the alternating pairs, which control for load drift, favour it, and the change is a soft limit (Go's collector still runs at the limit and caps its own CPU), so reverting it would discard a plausible small gain for nothing; keeping it is not counting it toward a target, since every 60% point is recorded as missed. The verifier's base re-run timing out at 600 s does not change this: the builder's base figures are consistent with the slice's base log (400 s wall for the suite, cmd/rhiza 254 s), and the commit side passed for both

## statement

Accept leaf cold-run-ship-scripts-cut at its commit 845741f94 (GOGC=off with GOMEMLIMIT=48MiB for every Go process a script starts, beside GOMAXPROCS=1), with no further attempt. Its 60% points are not met and not lowered: the three-package wall at 74.8-75.1% of the base's (verifier 371.79 s, builder 373.11 s against 496.88 s), cmd/rhiza's package time at 79.0-82.6% (against 231.001 s) and user+sys at 72.0-72.8% (against 857.92 s) are recorded as measured misses. The run-ship profile point (after's user+sys below before's) is also taken as not reliably met and not lowered: the acceptance's own pair shows no cut beyond noise on the loaded host (builder 89.60 -> 91.59 s; verifier 91.16 / 93.55, then 80.55 / 82.33), and the controlled alternating pairs show a small one (run-ship-conflict-fix 8.81 -> 8.27 s user+sys, user lower in all 4 pairs). Commit 845741f94 stays: it is test-harness only, in the leaf's files, skips no assertion (the scratch failure case still fails, 55 tests and 400 scripts unchanged), and no measurement shows it making anything worse beyond noise. The unbuilt -parallel default for a plain go test ./cmd/rhiza is not added, since make test already passes -parallel and it would not move the milestone's measurement. No trade-off against must_not is allowed: no internal/api edit, no non-test code change, no script line changed. The milestone is judged as m39-cold says, by the orchestrator's cold make test after the change set lands; internal/api's wall (about 370 s, the whole critical path) and the scripts' child-process cost in production code (git execs per command, rhiza process startup including internal/web init, sqlite index opens) are planned as change sets of their own if that measurement misses 240 s

## wording

(assumed, would have escalated)

## Links

### Out

- answers [cold-run-ship-scripts-cut-wall-unreachable](../question/cold-run-ship-scripts-cut-wall-unreachable.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
