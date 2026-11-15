# api-tests-parallel-session-qa-1

evidence · `01M7FJW538CC995K70ATJSFNJ9` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: The 214 newSessionFixture-file tests passed shuffled twice (seeds 7741, 90210; ~10 s each) and once under -race shuffled with no data race; ~/.cache/rhiza/run-tmp was unchanged across every run, and no new rhiza-api-* temp folder stayed after a normal run.

## Links

### Out

- verifies [api-tests-parallel-session](../capability/rhiza/integrity/self-check/api-tests-parallel-session.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
