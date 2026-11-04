# check-anchors-resolve

capability · `01M6HK9YY5HRN9EVKGV95EXMK1` · parent [anchors](../anchors.md)

## acceptance

check.Run stays pure: check.Input gains the working tree as an fs.FS rooted at the project, filled by internal/api's check; unit tests in internal/check with fstest.MapFS report one fail finding for a missing path and one for a missing symbol, naming the anchor and the capability it implements, set Report.OK false, and find nothing for a resolving anchor or a retired one; check-anchors.txtar: after deleting the anchored function, rhiza check exits non-zero with an anchor fail finding naming both, and after deleting the file a path fail finding; the rule is listed in the check package doc in report order

## statement

rhiza check has an anchor rule that resolves every live anchor of the folded model against the project's working tree and fails when its path or its symbol is gone, a broken anchor as the landed Decision ci-fails-on requires, each finding naming the anchor and the capability it implements

## Links

### Out

- depends-on [anchor-fields](anchor-fields.md)
- decided-by [anchor-fingerprint-shape](../../../../decision/anchor-fingerprint-shape.md)

### In

- depends-on [check-fingerprint-changed](../fingerprints/check-fingerprint-changed.md)
- verifies [check-anchors-resolve-verified](../../../../evidence/check-anchors-resolve-verified.md)
- verifies [check-anchors-resolve-qa-1](../../../../evidence/check-anchors-resolve-qa-1.md)
- verifies [check-anchors-resolve-verified-2](../../../../evidence/check-anchors-resolve-verified-2.md)
- verifies [pr-236-code-anchors-with-fingerprints](../../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
