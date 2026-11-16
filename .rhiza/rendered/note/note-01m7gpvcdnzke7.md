# note-01m7gpvcdnzke7

note · `01M7GPVCDNZKE783QC82Z46A7Y`

## statement

Index case (TestParityIndexTerminalVersion): one landed change set creating a node, write view of that landed change set. Before: index Version none (it left the change set's own ops out though it is not in flight); strict Version is the create. Index side changed: the current change set is the view's only while in flight.

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
