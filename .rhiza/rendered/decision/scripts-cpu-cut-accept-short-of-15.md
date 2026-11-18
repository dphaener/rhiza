# scripts-cpu-cut-accept-short-of-15

decision · `01M7PB903W2RN2MCANE33CSZSK`

## reasoning

This follows the slice's own precedent, api-tests-half-cpu-accept-at-60: a per-leaf CPU line is a proxy, not a milestone target; the intent's 'never lower a target' and m39-time-profile-first bind milestone targets, budget thresholds, rapid counts and tests, and accepting this leaf touches none of them. The builder profiled first as m39-time-profile-first requires, found the test process's own samples are only ~36 s of ~370 s (the rest in ~107k child processes, by strace), and cut the costs outside the scripts' assertions that the leaf's files allow; the verifier confirmed a real cut and every other acceptance point (tests pass, 55 tests and 400 scripts unchanged, the scratch failure case still fails). The two paths further both break the leaf's must_not: changing non-test code outside cmd/rhiza/main_test.go, or rewriting fake bin/ lines of scripts, which is deleting or changing lines of a script, so a retry cannot pass within the contract and failing the leaf would only block cold-run-ship-scripts-cut, which carries the wall-time checks (60% of the base's wall and cmd/rhiza package time, m39-cold-leaf-wall-proportion) that actually measure the milestone. GOGC and tmpfs were tried and did not help. Widening this leaf's files to non-test code would put production changes (git call counts, internal/web init) under a test-harness leaf without their own acceptance, so they belong in a separate change set if the downstream wall check shows they are needed.

## statement

Accept scripts-top-cost-cut at the cut commit 52037ffcb measured (GOMAXPROCS=1 for every Go process a script starts, git maintenance.auto and gc.auto off in the scripts' global config): a 9.2% to 12.6% user+sys cut of TestScripts (verifier 369.39 s to 335.55 s, builder 371.95 s to 324.96 s) is taken as meeting the leaf, its 15% line not met and not lowered. No trade-off against must_not is allowed: the leaf does not change non-test code (fewer git calls per command, a cheaper internal/web init) or rewrite the scripts' fake bin/ polling lines. Those costs, if needed to reach the milestone's 240 s, are new work in a change set of their own, not more attempts on this leaf; the milestone is judged by cold-run-ship-scripts-cut's wall-time check and the cold make test after landing

## wording

(assumed, would have escalated)

## Links

### Out

- answers [scripts-cpu-15pct-unreachable](../question/scripts-cpu-15pct-unreachable.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
