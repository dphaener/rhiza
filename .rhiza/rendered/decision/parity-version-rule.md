# parity-version-rule

decision · `01M7GGERQMXHJ6YYC2D5HYX4ZZ`

## question

Which version and conflict rule does the parity test hold the interim model and write_indexed.go to: the base graph's largest head with Conflict on two heads, or the strict engine's last applied operation with Conflict only on an overlapping stale write?

## reasoning

The owner made the strict engine the spec, and internal/model's doc already defines the version as the last operation applied, the one new writes name as base, and a conflict as a stale operation overlapping an aspect. The interim rule (head() in model.go, versionOf in write_indexed.go) sets Conflict on any two heads, so it refuses writes after two landed edits to different fields that the strict engine merges silently. Two heads writing the same field stay a conflict under both rules, so set.txtar and write_test.go's LandedConflict keep passing. Excluding the current change set's own operations matches what Version means on the api Node today. Pending is built the same way in both api models from the current change set's ops, so it is checked there.

## statement

The parity test compares a node's or link's Version and Conflict by the strict engine's rule: Version is the last operation in its strict history that is not the current change set's own (the version new writes name as base), and Conflict is true when the view holds an unrepaired model.ErrConflict problem on that target from an operation not of the current change set. Pending, which the strict engine has no notion of, is compared between the interim model and the index's write model only. A node or link one side has and the other lacks is a disagreement.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parity-harness-plan](../source/parity-harness-plan.md)

### In

- decided-by [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)
- decided-by [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)
- decided-by [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
