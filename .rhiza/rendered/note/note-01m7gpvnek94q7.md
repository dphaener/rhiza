# note-01m7gpvnek94q7

note · `01M7GPVNEK94Q7CWY1M60V9WPR`

## statement

Index case (TestParityIndexParentWithoutChildren): a draft creates a capability under a landed decision. Before: index holds it (present); strict skips it (model.ErrParent), absent. Index side changed.

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
