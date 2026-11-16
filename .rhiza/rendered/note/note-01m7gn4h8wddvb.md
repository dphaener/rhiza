# note-01m7gn4h8wddvb

note · `01M7GN4H8WDDVB788VE4V29MN9`

## statement

Case create-no-longer-applies (TestParityCreateNoLongerApplies): log = landed create of parent P; draft D creates child C under P; a later landed change set retires P. Write view of D: before, the interim model held C (present); the strict what-if skips D's create (model.ErrRetired) and has no C (absent). The interim side changed: newModelSets now takes its nodes and links from model.AsIsSeed and model.WhatIfFrom, so C is absent.

## Links

### Out

- discusses [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
