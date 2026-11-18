# cold-run-ship-scripts-cut

capability · `01M7P2X9HGS6ACMZCBMGJKGXSS` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza passes; go test -list '.*' ./cmd/rhiza and the count of cmd/rhiza/testdata/script/*.txtar are unchanged. Checks: block gives, before and after, /usr/bin/time -f '%e s %U user %S sys' go test -trimpath -count=1 -cpuprofile $TMPDIR/ship.prof ./cmd/rhiza -run 'TestScripts/run-ship' with the go tool pprof -top -cum lines of the cost cut, the after run's user+sys CPU below the before's; and /usr/bin/time -f '%e s %U user %S sys' go test -trimpath -count=1 ./internal/api ./cmd/rhiza ./internal/index run once at the commit and once in a git worktree at the change set's base 5d2596d under $TMPDIR, back to back: the commit's wall time at most 60% (240 s / 400 s) of the base's, cmd/rhiza's own package time (its go test ok line) at most 60% of the base's, and the commit's user+sys CPU at most 60% of the base's as a further check, both runs' wall, CPU and per-package times in the Checks: block. The full cold make test against m39-cold's 240 s is the milestone's own measurement, taken after the change set lands, not by the run. Failure case: a scratch copy of one run-ship script with its expected output changed still fails under the changed harness (named in the Checks: block, not committed), so the cut skips no assertion

## files

cmd/rhiza/main_test.go, cmd/rhiza/testdata/script/run-ship-*.txtar only where the cut moves a script's setup

## must_not

Add a script to processScripts; delete, skip or merge a script or a line of one; change non-test code outside cmd/rhiza/main_test.go; touch the network, Rex, agents or gh; regenerate cmd/rhiza/testdata/script-timings.tsv by hand; lower the 240 s target or any budget test threshold

## statement

The cold suite's critical path shortens: cmd/rhiza's run-ship scripts (run-ship-ci-fix, run-ship-conflict-fix, run-ship-conflict-fix-more, run-ship-change-rebase, run-ship-model-conflict, run-ship-land-conflict and run-ship-change, about 95 s of the 864 script-seconds in cmd/rhiza/testdata/script-timings.tsv) stop paying for the top cost a go test -cpuprofile of them names outside the commands they check after scripts-top-cost-cut, cut in cmd/rhiza's test harness; and once api-tests-copy-prebuilt-index and scripts-top-cost-cut have landed with it, a cold run of the three slowest packages is within m39-cold's 240 s target in proportion: its wall time, and cmd/rhiza's own package time, at most 240/400 of the base's

## Links

### Out

- depends-on [api-tests-copy-prebuilt-index](api-tests-copy-prebuilt-index.md)
- depends-on [scripts-top-cost-cut](scripts-top-cost-cut.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)
- decided-by [m39-cold-leaf-wall-proportion](../../../../decision/m39-cold-leaf-wall-proportion.md)

### In

- blocks [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f3](../../../../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f3.md)
- blocks [cold-run-ship-scripts-cut-wall-unreachable](../../../../question/cold-run-ship-scripts-cut-wall-unreachable.md)
- verifies [cold-run-ship-scripts-cut-verified-ruling](../../../../evidence/cold-run-ship-scripts-cut-verified-ruling.md)
- verifies [cold-run-ship-scripts-cut-qa-1](../../../../evidence/cold-run-ship-scripts-cut-qa-1.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
