# api-fixture-templates

decision · `01M7EQEPTN31N8JCPAY7A4ZC1B`

## reasoning

58% of internal/api's CPU is store.ReadLog because every write re-reads the log, so the cheap fix without projection is to make those writes once; a copied template keeps each test isolated (no test sees another's writes), which parallel tests need, where a shared read-only model would not let tests write.

## statement

internal/api's shared fixtures are built once per test binary into a template folder and copied into each test's own temporary folder, so a test still owns and may write its repo while the model-building writes run once

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [api-snapshot-fixtures](../capability/rhiza/integrity/self-check/api-snapshot-fixtures.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
