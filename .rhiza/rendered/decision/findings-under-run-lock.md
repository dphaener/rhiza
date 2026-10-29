# findings-under-run-lock

decision · `01M64JTD0AG6CRS3VBDC1J01R6`

## statement

Number a run's findings under the run's journal lock, reading the count and journalling in one locked section

## wording

Number findings under the run's lock, the agent's option (assumed, would have escalated)

## Links

### Out

- derived-from [findings-numbering-plan](../source/findings-numbering-plan.md)

### In

- decided-by [unique-finding-ids](../capability/unique-finding-ids.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
