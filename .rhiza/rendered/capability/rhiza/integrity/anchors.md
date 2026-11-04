# anchors

capability · `01M5SWT0Q3DHFNA33D1XP094VN` · parent [integrity](../integrity.md)

## statement

A capability is anchored to code by an anchor node, linked implements to it, that carries a path, an optional Go symbol and a content fingerprint, and rhiza check resolves every anchor against the working tree

## Links

### Out

- decided-by [anchor-fingerprint-shape](../../../decision/anchor-fingerprint-shape.md)

### In

- verifies [pr-236-code-anchors-with-fingerprints](../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
