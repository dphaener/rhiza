# scripts-parallel-bounded-qa-1

evidence · `01M7FPJNVYXXYA8ME3PT6DFTDN` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee: make -n test and make -n test-scripts SHARD=1/4 pass -parallel 8 by default on this 4-core VM and with NPROC=4 or 16, 4 with NPROC=2, 2 with NPROC=1, 0 or abc, 3 with RHIZA_TEST_SCRIPT_PARALLEL=3 (also as env, and over NPROC=16), and the NPROC default for 0, -2, abc and '3 4'; TestMakeTestTimeout and TestMakeTestShardTargets pass.

## Links

### Out

- verifies [scripts-parallel-bounded](../capability/rhiza/integrity/self-check/scripts-parallel-bounded.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
