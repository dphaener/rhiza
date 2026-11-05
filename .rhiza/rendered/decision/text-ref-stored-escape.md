# text-ref-stored-escape

decision · `01M6MNPXF7WMQ41TPW50H5Y1AS`

## reasoning

Only @@{<full 26-character ID>} counts as a stored reference (text-references-shape (2)), so storing @@@{<full ID>} as the literal @@{<full ID>} would read back as a reference nobody resolved: rhiza check fails it as dangling when the ID names no node, and it silently binds to the node when the ID names one, so a write that passed would leave the model wrong. A stored escape for it would need a second stored syntax and a change to every reader. Refusing in strict mode tells an author who can fix the text; the brace-less @<full ID> keeps an engine write from failing, which text-ref-engine-writes requires, keeps the ID readable, and is neither a stored nor a dangling reference. Recorded in review cycle 3, whose model check found that fix 22bd296a had changed these rules without a Decision.

## statement

Stored text never holds a literal @@{<full 26-character ID>}, since stored text cannot tell one from a reference, so the @@@ escape does not reach it. A strict write (rhiza add, rhiza set and the add and set items of rhiza apply) of @@@{<full ID>} is refused with a usage error (code usage, exit 2) naming it and its field, with a next saying to refer to the node with @@{<full ID>} or write the ID without braces, and writes nothing. A lenient engine write stores that escape, and any @@{<full ID>} that does not resolve, as @<full ID> without braces, never refusing; it is shown as @@@<full ID>. Every other @@@ still writes a literal @ that is never resolved. This narrows text-references-shape (2) ('@@ writes a literal @') and text-ref-engine-writes ('as if written @@@') for a full ID inside braces.

## wording

(assumed, would have escalated)

## Links

### In

- decided-by [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)
- decided-by [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
