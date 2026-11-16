# note-01m7gpvj2yy7w8

note · `01M7GPVJ2YY7W8ZBBM16TSTAGB`

## statement

Index case (TestParityIndexCreateTaken): a draft creates a node whose ID a landed create took. Before: index overwrote it (handle "other"); strict skips it (model.ErrExists, handle "store"). Index side changed.

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
