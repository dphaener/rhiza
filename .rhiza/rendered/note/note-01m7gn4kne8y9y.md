# note-01m7gn4kne8y9y

note · `01M7GN4KNE8Y9YTQZMSZDMCRGN`

## statement

Case disjoint-heads, the version rule of parity-version-rule (TestParityDisjointHeads): log = landed create of N; two landed updates both based on the create, one setting statement, the other acceptance. As-is: before, the interim model set Conflict true on two heads (version the larger head); the strict engine merges them, Conflict false, version the last op applied. The interim side changed. On this repo's own log it changes one node, 01M686TPQBCC7169KT6WEPB85A (question-answer-op: a node.move and a node.update both based on its create): Conflict goes from true to false in as-is and all 53 in-flight write views, Version unchanged, nothing else differs (compared old against new interim model over every in-flight view).

## Links

### Out

- discusses [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
