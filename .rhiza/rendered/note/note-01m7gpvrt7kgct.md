# note-01m7gpvrt7kgct

note · `01M7GPVRT7KGCT6TBT2WHN9H4G`

## statement

Index case (TestParityIndexDependencyApplyOrder): draft cur depends on a later draft dep (larger ID); both set statement of a landed node from its create. Before: index applied dep before cur, statement "cur", Conflict false; strict applies cur first (smallest ID), statement "dep", Conflict true (dep's stale write). Index side changed. The harness found no index case beyond these nine (the table is empty; the property test passes at the default and at -rapid.checks=1000, dry-run pending draws included).

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
