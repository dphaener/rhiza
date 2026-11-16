# note-01m7gpvq6601ma

note · `01M7GPVQ6601MA82GTEFE1KQCR`

## statement

Index case (TestParityIndexMoveCycle): a draft moves a landed node under its own child. Before: index parent is the child; strict skips the move (model.ErrCycle), parent none. Index side changed.

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
