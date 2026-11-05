# text-ref-syntax-qa-2

evidence · `01M6KVSR233SXSF0B0ABN4YC16` · actor [claude-code-e64add76](../actor/claude-code-e64add76.md)

## result

PASS

## statement

Scratch build of the branch head in a fresh repo: add capability with 'Builds on @@other (and [@@other]) mail darin@hey.com checkout@v7 @@@lit, prefix @<10-char ID prefix>' stored @@{<full ID>} for @@other, [@@other] and the prefix, kept darin@hey.com and checkout@v7 as written, stored @@@lit as @@lit; show printed '@other ... @@@lit'. A write of @@@{<26-char ID>} is now refused (exit 2, 'a literal @@{<full ID>} cannot be stored'), so the cycle-1 f2 hole is closed

## Links

### Out

- verifies [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
