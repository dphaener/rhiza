# api-tests-parallel-model-qa-1

evidence · `01M7FJW58KR5VMPQTESJHWMR4C` · actor [claude-code-89105d2a](../actor/claude-code-89105d2a.md)

## result

PASS

## statement

QA with test binaries built by go test -trimpath -c (and -race -c) to $TMPDIR, run from internal/api under the machine lock: The 201 model-group tests passed shuffled twice (seeds 31337, 424242; ~70 s each); under -race in two halves of 100/101 tests there was no data race; the only failure was TestRootIdlePollBudget's CPU budget under -race, which failed before this change set too.

## Links

### Out

- verifies [api-tests-parallel-model](../capability/rhiza/integrity/self-check/api-tests-parallel-model.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
