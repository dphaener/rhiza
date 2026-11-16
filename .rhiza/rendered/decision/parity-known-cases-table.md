# parity-known-cases-table

decision · `01M7GGERXZCDD3VHM2K6J2ZM8R`

## question

Does the parity test land red until the fixes land, or green with the disagreements it finds listed by name?

## reasoning

Every leaf's commit must pass CI before merge, and the fixes are a separate concern with their own regression tests. A named table keeps each case visible and makes the fix leaf's done point checkable (the table is empty), where skipped draws would hide cases.

## statement

The parity test lands green on its own leaf: a disagreement it finds goes into a named table of known cases that it reports and tolerates, and the next leaf fixes them and empties the table; the test never skips a draw to stay green.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parity-harness-plan](../source/parity-harness-plan.md)

### In

- decided-by [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
