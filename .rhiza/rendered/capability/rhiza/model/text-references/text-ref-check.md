# text-ref-check

capability · `01M6HP108NW407PT5CRA556JBA` · parent [text-references](../text-references.md)

## acceptance

Unit tests in internal/check: a stored @@{ID} naming no node gives one fail finding under a new rule text-ref, naming the node, the field and the ID, and sets Report.OK false; one naming a retired node gives a warn and keeps OK true; resolving references, text with no @@{, and literal look-alikes such as '@{ID}', '@{BOGUS}' and '@{<full ID>}' give nothing; check-text-refs.txtar: a log entry written by hand with a dangling @@{ID} makes rhiza check exit non-zero with that finding; the new rule appears in check.Rules() and in describe output; rhiza check on this repo stays clean, this plan's own statements and acceptances (which hold such look-alikes) included

## statement

rhiza check scans the folded model's stored references (only @@{<full 26-character ID>}; other @@{...} text is literal), each in the view of the change set that wrote it, and fails on one naming no node and warns on one naming a retired node

## Links

### Out

- depends-on [text-ref-syntax](text-ref-syntax.md)
- decided-by [text-references-shape](../../../../decision/text-references-shape.md)

### In

- verifies [text-ref-check-verified](../../../../evidence/text-ref-check-verified.md)
- verifies [text-ref-check-qa-1](../../../../evidence/text-ref-check-qa-1.md)
- verifies [text-ref-check-verified-2](../../../../evidence/text-ref-check-verified-2.md)
- verifies [text-ref-check-qa-2](../../../../evidence/text-ref-check-qa-2.md)
- verifies [text-ref-check-verified-7](../../../../evidence/text-ref-check-verified-7.md)
- verifies [pr-278-text-references](../../../../evidence/pr-278-text-references.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
