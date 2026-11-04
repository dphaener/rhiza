# anchor-fields

capability · `01M6HK9V1F5ZZNXE9RW11CZD1A` · parent [anchors](../anchors.md)

## acceptance

Unit tests of the resolver resolve a func, a method on a value and on a pointer receiver (Type.Method), a type, and a var and a const inside grouped declarations, hash exactly the declaration's source bytes (the span go/ast gives its node, doc comment excluded), and report a missing path, a missing symbol and a symbol on a non-.go file; anchor-fields.txtar: rhiza add anchor with path and symbol stores fingerprint sha256:<hex> of the declaration bytes, without symbol sha256:<hex> of the file (matching sha256sum), an explicit fingerprint field is kept as given, set of path or symbol recomputes it, and a missing path, a missing symbol and a symbol on a .md file are each refused with a Next; field-schema.txtar now expects path, symbol and fingerprint on an anchor; describe-golden.txtar is regenerated if a schema changed

## statement

An anchor carries path (relative to the project root, slash-separated), an optional symbol (a top-level Go declaration name, Type.Method for a method, accepted only on a .go file) and fingerprint (sha256:<hex> of the declaration's source bytes when symbol is set, else of the whole file's bytes, read from the working tree); adding an anchor, or setting its path or symbol, computes the fingerprint when none is given and refuses a path or symbol that does not resolve

## Links

### Out

- decided-by [anchor-fingerprint-shape](../../../../decision/anchor-fingerprint-shape.md)

### In

- depends-on [check-anchors-resolve](check-anchors-resolve.md)
- depends-on [anchor-refresh](../fingerprints/anchor-refresh.md)
- verifies [anchor-fields-verified](../../../../evidence/anchor-fields-verified.md)
- verifies [anchor-fields-qa-1](../../../../evidence/anchor-fields-qa-1.md)
- verifies [pr-236-code-anchors-with-fingerprints](../../../../evidence/pr-236-code-anchors-with-fingerprints.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
