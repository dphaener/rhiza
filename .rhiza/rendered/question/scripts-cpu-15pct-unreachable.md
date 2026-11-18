# scripts-cpu-15pct-unreachable

question · `01M7PAHHYYKK39VC5N5XE1QR98`

## statement

Leaf scripts-top-cost-cut: the 15% CPU cut (after user+sys at most 85% of the before) is out of reach inside the leaf's files and must_not. Measured back to back, go test -trimpath -count=1 -p 2 -cpuprofile ./cmd/rhiza -run TestScripts: before 371.95 s (196.85 user + 175.10 sys), after 324.96 s (168.71 + 156.25), a 12.6% cut (target <= 316.2 s); an earlier pair 364.27 -> 332.95 s with GOMAXPROCS alone. The cut: setup gives every Go process a script starts GOMAXPROCS=1 (each rhiza process's scheduler and GC threads spread over 4 cores), and the scripts' global git config turns off maintenance.auto and gc.auto (1645 detached git maintenance run --auto per suite). The test process's own pprof is only ~36 s of the ~370 s; the rest is ~107k child processes (strace): 44.5k git (rev-parse --git-common-dir 9k from the CLI's test-mode ID offset, user.email 2.9k, show-prefix 2.3k, rhiza.project.agent 2.1k, worktree list 1.9k, all from non-test rhiza code), ~6.1k rhiza processes (~11-13 ms CPU each to start, internal/web's init 4.6 ms of it), 4k sh, 3k sleep 0.025 polls in fake bin/ scripts. Tried and rejected: GOGC=200/400/off+GOMEMLIMIT (more CPU),  on tmpfs (no gain). Going further means changing non-test code (fewer git calls per command, a cheaper internal/web init) or rewriting the scripts' fake bin/ polling lines, both against must_not. Ruling asked: accept the leaf at 12.6%, or name which trade-off is allowed.

## Links

### Out

- blocks [scripts-top-cost-cut](../capability/rhiza/integrity/self-check/scripts-top-cost-cut.md)

### In

- answers [scripts-cpu-cut-accept-short-of-15](../decision/scripts-cpu-cut-accept-short-of-15.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
