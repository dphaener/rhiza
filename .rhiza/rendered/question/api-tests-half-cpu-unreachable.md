# api-tests-half-cpu-unreachable

question · `01M7P6J7WRJZ7G1F7FXXYGMKN9`

## statement

Leaf api-tests-copy-prebuilt-index: the half-CPU target (after user+sys at most 160.6 s, half of the before 321.25 s) is out of reach inside the leaf's files and must_not. Attempt 2 measured 193.69 s (144.08 user + 49.61 sys, 60% of before; attempt 1 199.77 s). Attempt 2 cut the 50,000-node model's throwaway index build in driverPollBudget (recordAda now records through ctx's index, TestReadBudget 39.9 s -> 29.9 s of samples) and moved TestRootIdlePoll's grown log into a snapshot (daemon-fold-98-grown-1100; the one build costs what the one catch-up did). What is left: GC 20 s; TestReadBudget 29.9 s, of which the 50,000-node model's one generation and index build (9.6 s) and the budgets' own setup writes, each new change set or actor record refolding that model's views (changeShowBudget change.new 4.8 s + warm-up 2.7 s, driverPollBudget newChange 2.5 s, actor record 4.1 s); TestRootIdlePoll 12.6 s (its 1100-change-set snapshot build); TestWritesThroughIndexMatchFold 11.4 s, whose fold side has no index on purpose (a throwaway index per write is its comparison); TestLinkEndpointKinds 6.3 s, TestPropertyProjectionParity 4.4 s. Cutting these means changing what the budgets measure on (setup writes baked into the generated model before it is indexed), sharing one index across runs of the binary, or dropping the fold side of TestWritesThroughIndexMatchFold, each against must_not or the test's meaning. Ruling asked: accept the leaf at 60% (the index setup the profile names cut as far as must_not allows), or name which of those trade-offs is allowed.

## Links

### Out

- blocks [api-tests-copy-prebuilt-index](../capability/rhiza/integrity/self-check/api-tests-copy-prebuilt-index.md)

### In

- answers [api-tests-half-cpu-accept-at-60](../decision/api-tests-half-cpu-accept-at-60.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
