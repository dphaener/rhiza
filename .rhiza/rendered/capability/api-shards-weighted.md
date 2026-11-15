# api-shards-weighted

capability · `01M7FT1545Y1MNAM7FRFQ2RZ1P`

## acceptance

A test on a fixture of timings (go test -trimpath -p 2 -count=1 ./tools/scripttimings or ./cmd/rhiza -run the new shard test) shows the shards disjoint and covering every listed test but the budgets, greedy longest first to the least-loaded shard, a missing test weighted at the median, ties to the lower shard; with the checked-in file the two shards' totals are within 15% of their mean. go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^TestMakeTestShardTargets$' passes, the union of make test-api's shards equal to go test -list minus BUDGET_TESTS; failure case: SHARD=0/2, 3/2 and 1/0 are still refused. Measure (orchestrator, on the next code PR): the two api jobs' Run api steps are within 20 s of each other and each under 120 s (2026-11-15, run 43721082089: 149 and 110 s). Touches Makefile, tools/scripttimings and cmd/rhiza/makeupdate_test.go. Must not change the CI matrix.

## statement

make test-api SHARD=i/n splits internal/api's top-level tests by recorded time instead of round-robin by name: each test, longest first by internal/api/testdata/test-timings.tsv (a missing one weighted at the median, for placement only), goes to the shard with the least total so far, ties to the lower shard, so CI's two api shards take about the same time.

## Links

### Out

- depends-on [api-test-timings](api-test-timings.md)
- decided-by [api-shards-weighted-like-scripts](../decision/api-shards-weighted-like-scripts.md)

### In

- verifies [api-shards-weighted-verified](../evidence/api-shards-weighted-verified.md)
- verifies [api-shards-weighted-qa-1](../evidence/api-shards-weighted-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
