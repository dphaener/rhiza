# text-ref-display-qa-1

evidence · `01M6KP0EYRFAJVFA3R0BD6JWR0` · actor [claude-code-5c694237](../actor/claude-code-5c694237.md)

## result

PASS

## statement

QA cycle 1 at 0e9f6a6e, scratch build in a scratch repo: after other was renamed to renamed, show, why, slice, search, change show and change diff printed @@renamed in text and --json; search '@renamed' found the node; after landing, rhiza render wrote @@renamed into capability/one.md; tui app show/why/slice one --once printed @@renamed. A retired reference showed @@{<full ID>} and writing it back with set stored the same text; writing one's shown statement back stored the same IDs and literal @. A hand-edited log entry 'Old @@word and @@{ID}' showed 'Old @@@word and @@@{ID}' and writing that back stored the same literal text; 'Lit @@@word' stored 'Lit @@word' and showed '@@word'. Verbatim quote fields and a source's summary print raw. (Hand-editing an entry needs a fresh XDG_CACHE_HOME: the existing index cache keys by entry.)

## Links

### Out

- verifies [text-ref-display](../capability/rhiza/model/text-references/text-ref-display.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
