# note-01m7gn4nj603gs

note · `01M7GN4NJ603GSVVVTEJKWFKQ0`

## statement

Beyond the named case and the version rule the harness found no other interim-side case: with the interim model built from the strict views, TestPropertyProjectionParity passes at the default checks and at -rapid.checks=1000 with every known case restricted to the index side. The two interim entries (Pending of an op before the target's create, and branches applied depth first in op ID order) are removed; the rest now name the index's write model, for parity-write-index-from-strict.

## Links

### Out

- discusses [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
