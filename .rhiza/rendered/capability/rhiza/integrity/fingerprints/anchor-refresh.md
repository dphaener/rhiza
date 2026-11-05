# anchor-refresh

capability · `01M6HKABFN71SDK6SXZAMPSXVF` · parent [fingerprints](../fingerprints.md)

## acceptance

anchor-refresh.txtar: after an anchored function changes, rhiza anchor refresh writes a set of the anchor's fingerprint into the current change set (--json and rhiza api call anchor.refresh print the same bytes), and once that lands rhiza check has no anchor finding; refreshing an unchanged anchor writes nothing and says so; refreshing an anchor whose path or symbol is gone is refused with a Next; --dry-run writes nothing; describe-golden.txtar is regenerated

## statement

rhiza anchor refresh <anchor>, one transport-independent operation in internal/api, re-fingerprints an existing anchor from the working tree into the current change set, the deliberate step after the code change is reviewed

## Links

### Out

- depends-on [anchor-fields](../anchors/anchor-fields.md)
- decided-by [anchor-fingerprint-shape](../../../../decision/anchor-fingerprint-shape.md)

### In

- verifies [anchor-refresh-verified](../../../../evidence/anchor-refresh-verified.md)
- verifies [anchor-refresh-qa-1](../../../../evidence/anchor-refresh-qa-1.md)
- verifies [pr-236-code-anchors-with-fingerprints](../../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
