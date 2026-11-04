# check-fingerprint-changed

capability · `01M6HKA5CCG3S9EABTJVXVBKHS` · parent [fingerprints](../fingerprints.md)

## acceptance

Unit tests in internal/check warn once for a changed declaration and once for a changed whole file, with Next rhiza anchor refresh <anchor>, find nothing when only bytes outside an anchored declaration change, and keep Report.OK true; check-fingerprint.txtar: after editing an anchored function's body, rhiza check exits 0 with one anchor warn finding naming the anchor and its capability, and make check would still pass

## statement

rhiza check's anchor rule warns, never fails, when an anchor's path and symbol resolve but the fingerprint of what they name in the working tree differs from the stored one, naming the anchor and the capability it implements, with Next the re-fingerprint command

## Links

### Out

- depends-on [check-anchors-resolve](../anchors/check-anchors-resolve.md)
- decided-by [anchor-fingerprint-shape](../../../../decision/anchor-fingerprint-shape.md)

### In

- verifies [check-fingerprint-changed-verified](../../../../evidence/check-fingerprint-changed-verified.md)
- verifies [check-fingerprint-changed-qa-1](../../../../evidence/check-fingerprint-changed-qa-1.md)
- verifies [pr-236-code-anchors-with-fingerprints](../../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
