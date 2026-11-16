# api-tests-on-index-qa-1

evidence · `01M7GVFM637RJ608A8G6WM3WG0` · actor [claude-code-4b0e81ca](../actor/claude-code-4b0e81ca.md)

## result

PASS

## statement

QA cycle 1: test-only leaf; ran the scoped parity/index tests above in a capped scope, ok; confirmed by hand that the CLI's production Ctx uses index.DefaultPath (cmd/rhiza/ctx.go:116-119), the path these tests now take.

## Links

### Out

- verifies [api-tests-on-index](../capability/rhiza/integrity/self-check/api-tests-on-index.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
