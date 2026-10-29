# pr-82-rhiza-run-findings-numbers

evidence · `01M64WKE783CNKRZ70WS5T1XMM` · actor [darin-haener](../actor/darin-haener.md)

## head

591dcd1402f407047e5f4bca53d5eb078a916a9b

## pr

https://github.com/dphaener/rhiza/pull/82

## result

1 check passed: check SUCCESS

## statement

Pull request #82 (https://github.com/dphaener/rhiza/pull/82) passed CI at 591dcd1

## status

verified

## Links

### Out

- verifies [unique-finding-ids](../capability/unique-finding-ids.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
