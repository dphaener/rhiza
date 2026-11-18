# cold-run-ship-scripts-cut-verified-ruling

evidence · `01M7PGDHH7PB23J7JTSKGZHAFA` · actor [claude-code-ae9599ba](../actor/claude-code-ae9599ba.md)

## result

PASS

## statement

Verified at 845741f94 per decision cold-run-ship-accept-wall-short-of-60 (no timed or whole-package runs, per the ruling). Commit's Checks: block holds the before/after run-ship pair, both pprof -top -cum lines, the per-PID child profile, and the per-package lines (internal/api 370.225s, cmd/rhiza 182.433s, internal/index 71.255s vs base), with the 60% points recorded as MISSED. go test -trimpath -p 2 -list '.*' ./cmd/rhiza: 55 Test lines at 52037ffcb and at 845741f94, identical (diff empty). txtar count 400 both sides, no testdata file in the diff. Failure case in a scratch worktree at 845741f94: run-ship-change.txtar:186 'landed' -> 'abandoned', go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/run-ship-change$' exit 1 with 'run-ship-change.txtar:186: no match for ^- State: abandoned found in stdout'; reverted, same command exit 0. Diff 52037ffcb..845741f94 touches only cmd/rhiza/main_test.go (GOGC=off, GOMEMLIMIT=48MiB in setup) plus .rhiza/log and the builder's rendered actor: no internal/api edit, no non-test code, no script line changed.

## Links

### Out

- verifies [cold-run-ship-scripts-cut](../capability/rhiza/integrity/self-check/cold-run-ship-scripts-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
