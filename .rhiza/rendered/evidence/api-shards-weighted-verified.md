# api-shards-weighted-verified

evidence · `01M7FXTNVHYXBKWCQK8CJTWCPD` · actor [claude-code-5fc8e63c](../actor/claude-code-5fc8e63c.md)

## result

PASS

## statement

At 1abcce176 I re-ran the commit's Checks and they all exit 0: go test -trimpath -p 2 -count=1 ./tools/scripttimings (TestShard covers greedy longest-first, a missing test at the median, ties to the lower shard, and disjoint/covering for n=1..4; TestShardCheckedInFileBalanced; TestShardMode refuses ''/0/2/3/2/1/0/1/2x/a/b/1), go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^TestMakeTestShardTargets$' (union of test-api's shards == go test -list minus BUDGET_TESTS; SHARD=0/2, 3/2 and 1/0 refused before go test runs), go vet ./tools/scripttimings ./cmd/rhiza, rhiza ship generate (0 changed) and rhiza render. I checked independently with the built tool on the 757 non-budget internal/api tests: shards 1/2 and 2/2 total 163.64 s each (260 and 497 tests, none missing from the TSV), with no duplicates and the union equal to the list. make test-api SHARD=0/2, 3/2 and 1/0 each exit 2 with 'want i/n with 1 <= i <= n'. The diff touches only Makefile, tools/scripttimings, cmd/rhiza/makeupdate_test.go and the builder's .rhiza records, has no .github/CI matrix change and no TODOs. The CI-timing measure is the orchestrator's, on the next code PR.

## Links

### Out

- verifies [api-shards-weighted](../capability/api-shards-weighted.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
