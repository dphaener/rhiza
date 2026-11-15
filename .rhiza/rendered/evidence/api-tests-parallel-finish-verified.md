# api-tests-parallel-finish-verified

evidence · `01M7FEJS87AADSQFTHK3QKAGYV` · actor [claude-code-ca930d0a](../actor/claude-code-ca930d0a.md)

## result

PASS

## statement

Verified at bf67c9599 (empty commit, only .rhiza records; no code or test changed). Under the machine lock with -p 2: go test -list '.*' ./internal/api counts 761 Test funcs (same as before, nothing skipped or deleted); go vet ./internal/api exit 0; whole-package go test -trimpath -count=1 with -skip '^(TestReadBudget|TestReadBudgetApply|TestDriverPollIsCheap)$' passed twice with new seeds -shuffle=91827 (158.5 s test, 160.5 s wall) and -shuffle=55301 (159.2 s, 160.5 s wall); about 10 s over the 150 s quiet-VM target on a busy shared VM, a note, not a failure. -race: built once with go test -race -c, run in sequential parts by file group, TestRootIdlePollBudget also skipped (pre-existing -race timing failure): daemon_* 52 tests 53.0 s, session_* 64 tests 5.1 s, ship_* 45 tests 9.7 s, run_* 113 tests 22.1 s, write_* 3 tests 74.9 s, other a-d 82 tests 38.8 s, and a new overlapping part daemon_|ship_|run_|c* 254 tests 99.7 s: every part exit 0 with 0 DATA RACE (7 of the builder's 10 parts re-run). rhiza ship generate and rhiza render exit 0, no tracked change beyond this verifier's own actor record.

## Links

### Out

- verifies [api-tests-parallel-finish](../capability/rhiza/integrity/self-check/api-tests-parallel-finish.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
