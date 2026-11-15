# test-scripts-sharded-greedy-verified

evidence · `01M7FN7Q0PTD4ENPD6Q1XP3G30` · actor [claude-code-0a2a9cf5](../actor/claude-code-0a2a9cf5.md)

## result

PASS

## statement

Verified at f8fa2173e: go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^(TestScriptShard|TestScriptShardEnvError|TestScriptOrder)$' (covers -run TestScriptShard) exit 0, every subtest PASS; gofmt -l cmd/rhiza empty; go vet ./cmd/rhiza exit 0. Diff read: shardNames sorts by recorded time (untimed weigh the median, timed first on ties, then name) and places each on the least-loaded shard, ties to the lower; scriptFiles passes the selection to testscript Files. TestScriptShard asserts: n=4 over testdata/script with script-timings.tsv disjoint and covering every script, each shard longest-first, totals within 15% of the mean; exact greedy placements pinned on fixtures (incl. tie to lower shard and median weight for missing scripts); unset returns all scripts longest first; 0/4, 5/4, x (and 1/0, 1/, /4, 1/4x, -1/4) rejected; TestScriptShardEnvError checks the error names RHIZA_TEST_SHARD with want i/n. Note: the n=4 real-data case checks balance rather than re-deriving the exact greedy assignment, which the fixtures pin. No TODOs or stubs.

## Links

### Out

- verifies [test-scripts-sharded](../capability/rhiza/integrity/self-check/test-scripts-sharded.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
