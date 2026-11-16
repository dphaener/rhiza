# note-01m7gpve6s3vax

note · `01M7GPVE6S3VAXPCP454WEPJ2F`

## statement

Index case, the version rule (TestParityIndexDisjointHeads): two landed updates on one create writing different fields. Before: index Conflict true (two heads); strict Conflict false, Version the last applied. Index side changed: Version is the last op of the strict history not the current change set's, Conflict an unrepaired model.ErrConflict from another change set (Decision parity-version-rule). On this repo's own log this is the only user-visible change: node question-answer-op (01M686TPQBCC7169KT6WEPB85A) no longer reports Conflict in as-is and the 53 in-flight write views, nothing else differs.

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
