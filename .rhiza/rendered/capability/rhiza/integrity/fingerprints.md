# fingerprints

capability · `01M5SWT0S5PN1WS8HSXEDY111W` · parent [integrity](../integrity.md)

## statement

A changed content fingerprint on an anchor makes rhiza check warn, naming the capability to review, until the anchor is deliberately re-fingerprinted

## Links

### Out

- decided-by [anchor-fingerprint-shape](../../../decision/anchor-fingerprint-shape.md)

### In

- verifies [pr-236-code-anchors-with-fingerprints](../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
