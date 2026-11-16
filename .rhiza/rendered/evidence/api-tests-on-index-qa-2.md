# api-tests-on-index-qa-2

evidence · `01M7GWB8FRTEBK9P9DCAGBMP98` · actor [claude-code-e084a010](../actor/claude-code-e084a010.md)

## result

PASS

## statement

QA cycle 2: go test ./internal/api -run '^(TestPropertyProjectionParity|TestParity|TestInboxUncoveredProposalsUnderLandedParent)' passed on the per-test index; the fix commit's TestInboxUncoveredProposalsUnderLandedParent compares fold and index reads and passes.

## Links

### Out

- verifies [api-tests-on-index](../capability/rhiza/integrity/self-check/api-tests-on-index.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
