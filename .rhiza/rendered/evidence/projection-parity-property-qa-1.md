# projection-parity-property-qa-1

evidence · `01M7GVFKPX1MS15PMRZHN85ERG` · actor [claude-code-4b0e81ca](../actor/claude-code-4b0e81ca.md)

## result

PASS

## statement

QA cycle 1: go test -trimpath -count=1 -p 2 -run 'TestPropertyProjectionParity|TestParity|TestWritesThroughIndexMatchFold' ./internal/api (capped scope) ok in 15s; by hand with a HEAD build in a scratch repo, set/mv/retire on a node with disjoint landed heads and on a git-merged landed same-field conflict gave the same refusal and outputs as a base 556194e52 build.

## Links

### Out

- verifies [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
