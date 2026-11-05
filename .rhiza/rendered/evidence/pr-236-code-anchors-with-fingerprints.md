# pr-236-code-anchors-with-fingerprints

evidence · `01M6JEN5A17R2QDADQVNNR4QKB` · actor [dphaener](../actor/dphaener.md)

## head

44d198720de15b2060e813165569184d4a001838

## pr

https://github.com/dphaener/rhiza/pull/236

## result

1 check passed: check SUCCESS

## statement

Pull request #236 (https://github.com/dphaener/rhiza/pull/236) passed CI at 44d1987

## status

verified

## Links

### Out

- verifies [anchors](../capability/rhiza/integrity/anchors.md)
- verifies [fingerprints](../capability/rhiza/integrity/fingerprints.md)
- verifies [anchor-fields](../capability/rhiza/integrity/anchors/anchor-fields.md)
- verifies [check-anchors-resolve](../capability/rhiza/integrity/anchors/check-anchors-resolve.md)
- verifies [check-fingerprint-changed](../capability/rhiza/integrity/fingerprints/check-fingerprint-changed.md)
- verifies [anchor-refresh](../capability/rhiza/integrity/fingerprints/anchor-refresh.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
