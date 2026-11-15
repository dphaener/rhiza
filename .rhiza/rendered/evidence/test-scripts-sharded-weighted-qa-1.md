# test-scripts-sharded-weighted-qa-1

evidence · `01M7FPJP4BNJEZJ3878N4MH8XN` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee: TestScriptShard passes; RHIZA_TEST_SHARD unset runs both named scripts, each of shards 1/4..4/4 runs a disjoint part of them, and malformed values (0/4, 5/4, x) fail TestScripts with a clear error.

## Links

### Out

- verifies [test-scripts-sharded](../capability/rhiza/integrity/self-check/test-scripts-sharded.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
