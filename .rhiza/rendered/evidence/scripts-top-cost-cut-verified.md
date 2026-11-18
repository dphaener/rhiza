# scripts-top-cost-cut-verified

evidence · `01M7PBB9NJKYXQKN22NT4EMAEH` · actor [claude-code-9c7ec9a0](../actor/claude-code-9c7ec9a0.md)

## result

PASS

## statement

Verified at 52037ffcb under ruling scripts-cpu-cut-accept-short-of-15 (CPU pair not rerun, per brief; earlier verifier pair 369.39 -> 335.55 s, builder 371.95 -> 324.96 s). The commit's Checks: block holds the before/after /usr/bin/time lines and the go tool pprof -top -cum lines of both profiles. Detached worktrees at eee179680 and 52037ffcb: go test -trimpath -p 2 -list '.*' ./cmd/rhiza exit 0 both, 55 tests, identical list; ls cmd/rhiza/testdata/script/*.txtar | wc -l = 400 both. Failure case: scratch edit of run-gate.txtar line 27 to stdout '^exit=7$', go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/^run-gate$' exit 1 with 'run-gate.txtar:27: no match for `^exit=7$` found in stdout'; reverted, the same command exit 0. git diff --stat eee179680 52037ffcb: only cmd/rhiza/main_test.go (GOMAXPROCS=1 in setup, maintenance.auto=false and gc.auto=0 in testGitConfig) plus .rhiza log and rendered actor; no non-test code, no txtar or fake bin/ line changed, no processScripts addition, so none of the ruling's forbidden trade-offs

## Links

### Out

- verifies [scripts-top-cost-cut](../capability/rhiza/integrity/self-check/scripts-top-cost-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
