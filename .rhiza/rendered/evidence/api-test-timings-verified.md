# api-test-timings-verified

evidence · `01M7FX95TPEE71YAEG8E00XGBH` · actor [claude-code-303d4398](../actor/claude-code-303d4398.md)

## result

PASS

## statement

At c63e4ff3f: go test -trimpath -p 2 -count=1 ./tools/scripttimings exit 0 (TestReplayTestTimes: test-time lines for passing top-level tests only, parsed back from log lines and from go test -json; TestMergeTests json and run-log: failed/unrun keep time, removed dropped, budget skipped, untimed new left out; TestMergeTestsFailuresLeaveFileUnchanged; TestCheckedInTestsFileParses). go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^TestMakeTestShardTargets$' exit 0 (fake go test: make test-api fails on failed test, failed test with exit 0, build failure, killed go test exit 137, and non-zero exit with passing stream; passes only when both pass). go test -trimpath -list '^Test' ./internal/api prints 762 tests; minus the 5 BUDGET_TESTS present there, the 757 names equal internal/api/testdata/test-timings.tsv's names exactly (diff empty). go vet and gofmt clean on tools/scripttimings; no TODOs; changes stay within Makefile, tools/scripttimings, the new tsv, cmd/rhiza/makeupdate_test.go and .rhiza records.

## Links

### Out

- verifies [api-test-timings](../capability/api-test-timings.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
