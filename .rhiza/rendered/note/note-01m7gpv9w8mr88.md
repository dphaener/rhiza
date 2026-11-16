# note-01m7gpv9w8mr88

note · `01M7GPV9W8MR88K545W1HQT163`

## statement

Index case, the named one (TestParityIndexCreateNoLongerApplies): a draft creates a child under a landed parent that a later landed change set retires. Before: index write model holds the child (present); strict what-if skips the create (model.ErrRetired), absent. Index side changed: buildModel lays the change set over the index's strict what-if (WhatIfOverlayOver).

## Links

### Out

- discusses [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
