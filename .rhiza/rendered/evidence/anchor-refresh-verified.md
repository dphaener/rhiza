# anchor-refresh-verified

evidence · `01M6J7DPH0Z9S9312922YCGW77` · actor [claude-code-746198af](../actor/claude-code-746198af.md)

## result

PASS

## statement

At 02556758, go test ./cmd/rhiza -run 'TestScripts/(anchor-refresh|describe-golden)$' -count=1 passes. anchor-refresh.txtar covers each acceptance point: after the anchored function changes, anchor refresh writes a set of the fingerprint into the current change set; --json and api call anchor.refresh print the same bytes (cmp, on a dry run); after land plus render, check is ok with no anchor finding and a further refresh says unchanged; refreshing an unchanged anchor writes nothing and says so; a gone path, a gone symbol or a non-anchor is refused with a next line; --dry-run writes nothing. describe-golden passes with the new op's schema. I read the diff: internal/api/anchor_refresh.go registers one write op that reuses runWrite and buildSet, and an unchanged anchor returns before any entry is appended. A manual run in a scratch repo with a fresh build agreed: refresh wrote previous to new fingerprint, then api call reported unchanged, and check stayed ok after land. make check: fmt-check, vet, self-check and every test pass except three read-timing budgets (TestReadBudget, TestBudget, TestGeneratedModel). Those failed under load average about 20 with other make checks running beside them, and they pass with only the timings skipped. The diff touches no read path.

## Links

### Out

- verifies [anchor-refresh](../capability/rhiza/integrity/fingerprints/anchor-refresh.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
