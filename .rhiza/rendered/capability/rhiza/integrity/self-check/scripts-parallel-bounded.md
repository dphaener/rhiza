# scripts-parallel-bounded

capability · `01M7EQE64JDCYTPPV1369XDSNM` · parent [self-check](../self-check.md)

## acceptance

make -n test and make -n test-scripts SHARD=1/4 pass -parallel 8 with NPROC=4 (the Makefile's CPU count: getconf _NPROCESSORS_ONLN unless NPROC is set), -parallel 4 with NPROC=2, -parallel 8 with NPROC=16, and -parallel 3 with RHIZA_TEST_SCRIPT_PARALLEL=3; go test -trimpath ./cmd/rhiza -run '^(TestMakeTestTimeout|TestMakeTestShardTargets)$' covers each; failure case: RHIZA_TEST_SCRIPT_PARALLEL=0, -2 or abc falls back to the NPROC default rather than passing -parallel 0 or failing make. Speedup and memory measure (orchestrator, quiet 4-core VM), through the Makefile so its -parallel applies: go clean -testcache && systemd-run --user --scope -p MemoryMax=6G -p TasksMax=300 -- /usr/bin/time -f '%e s %M KB' make test-scripts passes (at the default -parallel 8) without hitting either cap, and takes less wall time than the same command after go clean -testcache with RHIZA_TEST_SCRIPT_PARALLEL=4 make test-scripts. Touches Makefile, cmd/rhiza/makeupdate_test.go. Must not raise the default's cap above 8 or drop the machine-wide test lock.

## statement

make test and make test-scripts run cmd/rhiza's scripts at a -parallel of twice the CPU count capped at 8, or RHIZA_TEST_SCRIPT_PARALLEL when it is a positive integer, instead of go test's default of GOMAXPROCS.

## Links

### Out

- decided-by [script-parallel-bound](../../../../decision/script-parallel-bound.md)

### In

- verifies [scripts-parallel-bounded-verified](../../../../evidence/scripts-parallel-bounded-verified.md)
- verifies [scripts-parallel-bounded-qa-1](../../../../evidence/scripts-parallel-bounded-qa-1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
