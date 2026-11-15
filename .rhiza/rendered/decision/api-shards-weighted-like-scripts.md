# api-shards-weighted-like-scripts

decision · `01M7FT150122VYYP8MBEW2QPFR`

## reasoning

The brief asks to rebalance the api CI shards. Round-robin by name put the 81 s TestLinkEndpointKinds and half the serial tests on one shard (149-165 s against 110-137 s); a timing file keeps the split even as tests are added. Reusing scripttimings and #593's greedy split keeps one mechanism. Two shards, not more, because each extra job raises the chance of waiting for a runner.

## statement

CI's internal/api shards split by recorded test time the way the scripts shards do: make test-api pipes go test -json through tools/scripttimings replay, which prints a test-time line per passing top-level test; make api-timings RUN=<id> or JSON=<file> merges them into a checked-in internal/api/testdata/test-timings.tsv (seeded from one local go test -json run of internal/api, budgets skipped); and test-api gives each top-level test, longest first, to the shard with the least total so far, a test missing from the file weighted at the median, ties to the lower shard. Still two shards.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [api-test-timings](../capability/api-test-timings.md)
- decided-by [api-shards-weighted](../capability/api-shards-weighted.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
