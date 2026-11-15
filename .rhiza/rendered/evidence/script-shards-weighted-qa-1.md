# script-shards-weighted-qa-1

evidence · `01M7FPJP0JWWTSMV59S7VWTCFG` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee: TestScriptShard and TestScriptShardEnvError pass; real TestScripts with RHIZA_TEST_SHARD=0/4, 5/4, 1/0, x and 1/4x fails with 'want i/n with 1 <= i <= n'; running version and describe-golden under shards 1/4..4/4 runs each in exactly one shard (describe-golden in 2, version in 3).

## Links

### Out

- verifies [script-shards-weighted](../capability/rhiza/integrity/self-check/script-shards-weighted.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
