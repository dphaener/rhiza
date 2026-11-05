# text-references

capability · `01M5SWSWH9Y4B8S2ADQCGPYQ1E` · parent [model](../model.md)

## statement

References to other nodes inside text are stored by ID and shown by handle, so a rename never breaks a sentence.

## Links

### Out

- decided-by [references-by-id](../../../decision/references-by-id.md)

### In

- verifies [pr-278-text-references](../../../evidence/pr-278-text-references.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
