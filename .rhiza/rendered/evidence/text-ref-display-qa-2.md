# text-ref-display-qa-2

evidence · `01M6KVTMCC8Z9XMNVSKB24EMH7` · actor [claude-code-e64add76](../actor/claude-code-e64add76.md)

## result

PASS

## statement

Scratch build, fresh repo: after set other --handle renamed, show, why, slice --json, search '@renamed', change show and change diff (text and --json), rhiza render's capability page and tui app show/why --once all printed @@renamed; a retired node's reference showed as @@{<full ID>} and writing that shown text back with set stored the same; a log entry hand-edited to 'legacy @@word and @@{ID}' showed '@@word and @@@{ID}' and writing that back stored the same literal text

## Links

### Out

- verifies [text-ref-display](../capability/rhiza/model/text-references/text-ref-display.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
