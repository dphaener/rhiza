# rebase-conflict-text-references

question · `01M6MSPG3Y1KBCTD94GGCSPTNE`

## statement

Rebasing branch cs/text-references onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: .rhiza/rendered/index.md, cmd/rhiza/testdata/script/describe-golden.txtar, internal/api/check.go, internal/check/check.go, internal/check/check_test.go?

## Links

### Out

- blocks change set `01M6HNNWFR1J2AHQRNKFTPME46` (Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model))

### In

- answers [rebase-conflict-text-references-resolved](../decision/rebase-conflict-text-references-resolved.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
