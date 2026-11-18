# cold-run-ship-scripts-cut-wall-unreachable

question · `01M7PEV3MRAQSJ4SDNDPJRHEKR`

## statement

Leaf cold-run-ship-scripts-cut, attempt 2: two of its wall-time points are out of reach inside the leaf's files (cmd/rhiza/main_test.go, run-ship scripts) and must_not. (1) Three-package wall at most 60% of base: attempt 1 measured 373.11 s against base 496.88 s (75.1%); internal/api alone is 370 s of it, outside this leaf's files, and the brief forbids editing internal/api here. Reaching 60% needs internal/api at about 298 s or less. (2) cmd/rhiza's own package time at most 60% of base (231.001 s, so at most 138.6 s): attempt 1 182.4 s (verifier 186.5 s, 79.0%). Attempt 2 measured go test -trimpath -count=1 -p 2 -parallel 8 ./cmd/rhiza (the -parallel make test already passes) at ok 161.516 s, 162.68 s wall, 203.28 user + 182.68 sys = 385.96 CPU-s (load 4 at start, 9.7 at end from other sessions): the box is CPU-bound at about 2.4 of 4 cores, and even perfect packing of that CPU on 4 cores is about 97 s. The CPU is in child processes, not the test process (earlier strace: about 36 of about 370 CPU-s in the test process; about 107k child processes, 44k of them git and 6k rhiza). Attempt 2's strace of run-ship-ci-fix: 1872 execs, 383 fsyncs costing 0.34 s of a 16 s script, so there are no I/O waits for the harness to remove. Harness levers already used: GOMAXPROCS=1, git auto-maintenance and gc off, GOGC=off with GOMEMLIMIT=48MiB. tmpfs and GOGC=off at 128MiB were tried and did not help. What is left is production code: git calls per command, about 14 ms of startup per rhiza process (internal/web init) and sqlite index opens. Ruling scripts-cpu-cut-accept-short-of-15 already put those outside a test-harness leaf. (3) user+sys at most 60% of base: 617.71 s against 857.92 s (72.0%), with internal/api's CPU in it. One harness change could move the cmd/rhiza number: making a plain go test ./cmd/rhiza run scripts at -parallel min(2*NPROC,8) when -parallel is not given, as make test does. It still misses (161.5 s against 138.6 s). It also would not shorten the cold make test the milestone measures, since make test already passes -parallel 8, so I have not made it. Ruling asked: accept the leaf at attempt 1's measured cut, its 60% wall points not met and not lowered, with internal/api's wall and the scripts' child-process CPU (git calls, rhiza startup) planned as their own change sets; or name a cut inside this leaf's files and must_not to build.

## Links

### Out

- blocks [cold-run-ship-scripts-cut](../capability/rhiza/integrity/self-check/cold-run-ship-scripts-cut.md)

### In

- answers [cold-run-ship-accept-wall-short-of-60](../decision/cold-run-ship-accept-wall-short-of-60.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
