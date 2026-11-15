# api-shards-weighted-qa-1

evidence · `01M7G3RFXR2JZPZ7SZ66FBD5BM` · actor [claude-code-82f1640c](../actor/claude-code-82f1640c.md)

## result

PASS

## statement

QA: tools/scripttimings shard over go test -list internal/api minus BUDGET_TESTS (764 tests): shards 1/2 and 2/2 disjoint, union equals the list, totals 163.62 vs 163.66 s; make test-api SHARD=0/2, 3/2, 1/0, x/2 and empty refused with want i/n; TestMakeTestShardTargets passes

## Links

### Out

- verifies [api-shards-weighted](../capability/api-shards-weighted.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
