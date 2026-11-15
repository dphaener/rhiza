# api-tests-parallel-session-verified

evidence · `01M7F3WHK81G5DEZSAPKBY1HP5` · actor [claude-code-17fef95a](../actor/claude-code-17fef95a.md)

## result

PASS

## statement

Verified at 1238b49e0 (flock + capped scope, -p 2, session env unset): the acceptance's -run pattern (214 tests in the files calling newSessionFixture) passed with -shuffle=on seed 1794729880646191513 and again with seed 12345 (214 PASS each, exit 0); the same run with -race exited 0 with no DATA RACE; go vet ./internal/api exit 0; go test -list '.*' ./internal/api lists 756 Test funcs, and git grep counts 756 top-level Test funcs at both HEAD^ and HEAD (the diff removes only the fixture's t.Setenv lines, no test). newSessionFixture has no t.Setenv; TestMain in budget_test.go MkdirTemps rhiza-api-cache-, sets XDG_CACHE_HOME, sweeps stale ones and RemoveAlls it before exit (no rhiza-api-cache-* left in /tmp, ~/.cache/rhiza/run-tmp unchanged at 3 entries). 21 tests stay serial, each with a 'Serial:' comment (runClock/beforeDoneMark swap, t.Setenv, run-tmp folders); TestAdvisorBrief delegates to the parallel TestArbiterBrief, commented. Only _test.go files and .rhiza actor/log/render files changed.

## Links

### Out

- verifies [api-tests-parallel-session](../capability/rhiza/integrity/self-check/api-tests-parallel-session.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
