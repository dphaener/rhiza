# pr-278-text-references

evidence · `01M6N0R4GY23KT38PAYZANK90Z` · actor [dphaener](../actor/dphaener.md)

## head

7d01b3d176a8003f60f1cf221c7b275ebebdb73b

## pr

https://github.com/dphaener/rhiza/pull/278

## result

1 check passed: check SUCCESS

## statement

Pull request #278 (https://github.com/dphaener/rhiza/pull/278) passed CI at 7d01b3d

## status

verified

## Links

### Out

- verifies [text-references](../capability/rhiza/model/text-references.md)
- verifies [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)
- verifies [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)
- verifies [text-ref-display](../capability/rhiza/model/text-references/text-ref-display.md)
- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
