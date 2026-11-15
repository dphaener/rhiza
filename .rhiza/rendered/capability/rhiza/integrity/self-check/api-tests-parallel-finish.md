# api-tests-parallel-finish

capability · `01M7F32H4Y94GASPX3X98497EV` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 -shuffle=on ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' passes twice in a row with different shuffle seeds; go vet ./internal/api is clean. Failure case: go test -trimpath -race -count=1 ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' reports no data race; a cross-group race or shuffle failure is fixed by making that test serial with a comment saying why, or by fixing its shared state. Speedup measure (orchestrator, quiet VM): /usr/bin/time -f '%e s' go test -trimpath -count=1 ./internal/api -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' within 150 s (2026-11-15: 388 s). May touch any internal/api _test.go file. Must not skip, delete or merge a test: go test -list '.*' ./internal/api counts the same tests as before. Must not change non-test code.

## statement

With every internal/api test group parallelised and the snapshot fixtures in, the whole internal/api package passes shuffled, under -race and within the 150 s wall-time target, any cross-group failure fixed in whichever _test.go file it shows up in.

## Links

### Out

- depends-on [api-tests-parallel-session](api-tests-parallel-session.md)
- depends-on [api-tests-parallel-runtime](api-tests-parallel-runtime.md)
- depends-on [api-tests-parallel-rest](api-tests-parallel-rest.md)
- depends-on [api-tests-parallel-model](api-tests-parallel-model.md)
- depends-on [api-snapshot-fixtures](api-snapshot-fixtures.md)
- decided-by [api-parallel-rest-split-and-finish](../../../../decision/api-parallel-rest-split-and-finish.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [api-cache-per-binary](../../../../decision/api-cache-per-binary.md)

### In

- verifies [api-tests-parallel-finish-verified](../../../../evidence/api-tests-parallel-finish-verified.md)
- verifies [api-tests-parallel-finish-qa-1](../../../../evidence/api-tests-parallel-finish-qa-1.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
