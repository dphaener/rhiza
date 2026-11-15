# scripts-parallel-bounded-verified

evidence · `01M7FNS63S1V62VRTKYY1VFZXW` · actor [claude-code-81e6dc31](../actor/claude-code-81e6dc31.md)

## result

PASS

## statement

At 6e4fed7ee: go test -trimpath -p 2 ./cmd/rhiza -run '^(TestMakeTestTimeout|TestMakeTestShardTargets)$' -count=1 exit 0 (both PASS; checkScriptParallel covers NPROC 4/2/16/1, overrides 3/12, and 0/-2/abc/2x fallbacks for make test and make test-scripts SHARD=1/4); rhiza ship generate exit 0 (0 files changed); rhiza render exit 0. By hand, make -n for both make test and make test-scripts SHARD=1/4 printed -parallel 8 (NPROC=4), 4 (NPROC=2), 8 (NPROC=16), 3 (RHIZA_TEST_SCRIPT_PARALLEL=3), and the NPROC default 8/4/8 for RHIZA_TEST_SCRIPT_PARALLEL=0/-2/abc, all exit 0; the env-var form works too; the default on this 4-core VM is -parallel 8. -parallel appears once, only on cmd/rhiza's go test (test-rest and test-budgets have none); the default is capped at 8; lock "the cmd/rhiza scripts" is still taken in both targets. Diff touches only Makefile, cmd/rhiza/makeupdate_test.go and the builder's actor record; no TODOs. The speedup/memory measure was left to the orchestrator as briefed.

## Links

### Out

- verifies [scripts-parallel-bounded](../capability/rhiza/integrity/self-check/scripts-parallel-bounded.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
