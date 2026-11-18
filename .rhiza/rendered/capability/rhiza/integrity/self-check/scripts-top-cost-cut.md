# scripts-top-cost-cut

capability · `01M7P20X7FRK4EP2SSK5VYVRPZ` · parent [self-check](../self-check.md)

## acceptance

Checks: block gives, before and after, /usr/bin/time -f '%e s %U user %S sys' go test -trimpath -count=1 -cpuprofile $TMPDIR/scripts.prof ./cmd/rhiza -run TestScripts and the go tool pprof -top -cum lines of the cost cut, the after run's user+sys CPU at least 15% below the before's; go test -trimpath -count=1 ./cmd/rhiza passes; go test -list '.*' ./cmd/rhiza and the count of cmd/rhiza/testdata/script/*.txtar are unchanged. Failure case: a scratch copy of one affected script with its expected output changed still fails under the changed harness (named in the Checks: block, not committed), so the cut skips no assertion

## files

cmd/rhiza/main_test.go, cmd/rhiza/testdata/script/*.txtar only where the cut moves a script's setup

## must_not

Add a script to processScripts; delete, skip or merge a script or a line of one; change non-test code outside cmd/rhiza/main_test.go; touch the network, Rex, agents or gh; regenerate cmd/rhiza/testdata/script-timings.tsv by hand

## statement

cmd/rhiza's testscripts stop paying for the top cost a go test -cpuprofile of ./cmd/rhiza names outside the commands they check (setup in main_test.go's Setup, per-script repository or index building, git or process starts the scripts do not assert on), cut in the test harness so every script still runs every line

## Links

### Out

- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [cold-run-ship-scripts-cut](cold-run-ship-scripts-cut.md)
- blocks [scripts-cpu-15pct-unreachable](../../../../question/scripts-cpu-15pct-unreachable.md)
- verifies [scripts-top-cost-cut-verified](../../../../evidence/scripts-top-cost-cut-verified.md)
- verifies [scripts-top-cost-cut-qa-1](../../../../evidence/scripts-top-cost-cut-qa-1.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
