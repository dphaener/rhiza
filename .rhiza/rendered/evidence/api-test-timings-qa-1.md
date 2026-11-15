# api-test-timings-qa-1

evidence · `01M7G3RG00QH4D3MZZVEWQH3HV` · actor [claude-code-82f1640c](../actor/claude-code-82f1640c.md)

## result

PASS

## statement

QA: go test ./tools/scripttimings passes; internal/api/testdata/test-timings.tsv lists 757 of 764 non-budget tests (7 added by this change set missing, weighted at median by shard); TestMakeTestShardTargets pins test-api exit status

## Links

### Out

- verifies [api-test-timings](../capability/api-test-timings.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
