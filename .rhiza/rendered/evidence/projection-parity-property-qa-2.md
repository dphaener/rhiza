# projection-parity-property-qa-2

evidence · `01M7GWB88RW2QZBEY1K357V60E` · actor [claude-code-e084a010](../actor/claude-code-e084a010.md)

## result

PASS

## statement

QA cycle 2: go test ./internal/api -run '^(TestPropertyProjectionParity|TestParity)' -v passed: 100 draws, 52 with two landed heads, 75 with a rebuild, 35 with pending entries; TestParityReportNamesNode passed.

## Links

### Out

- verifies [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
