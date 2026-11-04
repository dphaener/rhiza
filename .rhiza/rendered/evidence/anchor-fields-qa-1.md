# anchor-fields-qa-1

evidence · `01M6JAE2EMCFY62SRTK5M25MVJ` · actor [claude-code-666954d5](../actor/claude-code-666954d5.md)

## result

PASS

## statement

Built the branch head (9a4dfa5d) into a scratch binary and drove it by hand in a scratch git repo with RHIZA_ROOT/CHANGE/RUN unset. add anchor with symbol=Keep (func), T.M (pointer-receiver method) and B (const in a grouped block) stored sha256 values that match sha256sum of exactly the declaration bytes, doc comment excluded; without symbol, notes.md got sha256sum of the file; an explicit fingerprint=sha256:abc was kept as given; set cb --field symbol=A recomputed it to sha256 of 'A = 1'. A missing path, a missing symbol and a symbol on a .md file were each refused with exit 2 and a next line.

## Links

### Out

- verifies [anchor-fields](../capability/rhiza/integrity/anchors/anchor-fields.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
