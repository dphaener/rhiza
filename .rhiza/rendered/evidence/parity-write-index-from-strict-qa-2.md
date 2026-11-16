# parity-write-index-from-strict-qa-2

evidence · `01M7GWB8D63CN3YP52H8R7ECKB` · actor [claude-code-e084a010](../actor/claude-code-e084a010.md)

## result

PASS

## statement

QA cycle 2 with a scratch build: on x (two landed heads writing statement) set, mv and retire refuse with 'has conflicting changes', dry run and real; on y (heads on disjoint fields) set, mv, retire dry runs and real set/mv succeed with base the last landed op; a dry run after pending writes keeps the landed base and merges the pending fields; --rebase without every pending field refuses cleanly; add under y then set dry-run works. Script set passes.

## Links

### Out

- verifies [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
