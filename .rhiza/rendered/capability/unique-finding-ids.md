# unique-finding-ids

capability · `01M64JTDQ1PN9MMS77P95QM4EG`

## acceptance

A testscript runs concurrent findings calls for one round, one --check and one leaf, then asserts that the findings journalled within each of those scopes have unique IDs

## statement

Concurrent findings calls within one round, check or leaf never journal two findings with the same ID

## Links

### Out

- decided-by [findings-under-run-lock](../decision/findings-under-run-lock.md)

### In

- verifies [unique-finding-ids-verified](../evidence/unique-finding-ids-verified.md)
- verifies [unique-finding-ids-qa-1](../evidence/unique-finding-ids-qa-1.md)
- verifies [pr-82-rhiza-run-findings-numbers](../evidence/pr-82-rhiza-run-findings-numbers.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
