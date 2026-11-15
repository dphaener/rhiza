# api-test-timings

capability · `01M7FT1523ZQDRWQFA9F4JW9S3`

## acceptance

go test -trimpath -p 2 -count=1 ./tools/scripttimings passes, with tests for the test-time lines (printed for top-level tests only, parsed back from a gh run log line and from go test -json) and for the api merge (a failed or unrun test keeps its time, a removed test is dropped). go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^TestMakeTestShardTargets$' pins that make test-api keeps go test's exit status through the pipe (fails on a failed test, a failed build or a killed go test) as test-scripts does. The checked-in file lists every top-level test go test -list '^Test' ./internal/api prints except the BUDGET_TESTS ones, from one go test -trimpath -p 2 -json ./internal/api -skip '<BUDGET_TESTS>' run. Touches Makefile, tools/scripttimings/*.go and tests, internal/api/testdata/test-timings.tsv (new) and cmd/rhiza/makeupdate_test.go.

## statement

Each CI api job's log carries its tests' times and a checked-in file holds them: make test-api pipes go test -json through tools/scripttimings replay, which prints go test's text for what did not pass and one test-time line per passing top-level test, and make api-timings RUN=<id> or JSON=<file> merges those times into internal/api/testdata/test-timings.tsv, seeded from a local go test -json run of internal/api.

## Links

### Out

- decided-by [api-shards-weighted-like-scripts](../decision/api-shards-weighted-like-scripts.md)

### In

- depends-on [api-shards-weighted](api-shards-weighted.md)
- blocks [review-01m7frv7b35b91y35xa8vxzpc0-r3-f1](../question/review-01m7frv7b35b91y35xa8vxzpc0-r3-f1.md)
- verifies [api-test-timings-verified](../evidence/api-test-timings-verified.md)
- verifies [api-test-timings-qa-1](../evidence/api-test-timings-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
