# api-tests-parallel-model-verified

evidence · `01M7F97XNTT8TSWA8R1TTKZ89K` · actor [claude-code-e48d9bb0](../actor/claude-code-e48d9bb0.md)

## result

PASS

## statement

At 827a1be3d (all with -p 2, under the machine lock): go vet ./internal/api exit 0. The acceptance -run pattern (196 tests from the 52 files) with the acceptance -skip, -count=1: -shuffle=101 ok (99.8 s), -shuffle=20261115 ok (108.5 s). -race split by file (files 1-26 / 27-52 of the acceptance list): first half with TestRootIdlePollBudget also skipped exit 0, 0 DATA RACE (108 s); second half with TestLinkEndpointKinds also skipped exit 0, 0 DATA RACE (116 s). TestRootIdlePollBudget alone with -race fails its CPU budget (~175 ms vs 100 ms) at 827a1be3d and the same way at its parent 64682feb7 (exit 1 both), so it is a pre-existing -race-only failure, not this leaf's. TestLinkEndpointKinds was not run with -race (73 s without -race per the builder; would exceed the 4-minute cap); it passed in both shuffled runs. go test -list '.*' ./internal/api: 768 entries, identical sorted list at HEAD and HEAD^ (nothing skipped, deleted or merged). Diff read: only internal/api _test.go files (plus the .rhiza actor record) change; 183 added t.Parallel calls and 13 '// Serial:' comments. Each serial reason checked: t.Chdir (TestInit), t.Setenv via exposureRepo/sliceCtx (every caller in these files is serial), the package merges counter (TestDaemonModelPR), process-wide read/fold counters (TestRootIdlePoll), timing (TestReadBudget, TestReadBudgetApply, TestRootIdlePollBudget). A scan of every package var written in these files found no parallel test that writes one (diskFreeBytes/machineReading only in TestMain, prePushCheck only in init). No stubs or TODOs.

## Links

### Out

- verifies [api-tests-parallel-model](../capability/rhiza/integrity/self-check/api-tests-parallel-model.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
