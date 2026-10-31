# daemon-skips-saved-at-poll

decision · `01M686SD9VM3411J6KDG2ZNPVK`

## statement

The daemon works out its skip reasons at each poll, including waiting for a slot, and saves them in its state file; daemon status reads them back as of the last poll

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m7-keep-picking](../source/m7-keep-picking.md)

### In

- decided-by [daemon-skip-reasons](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons.md)
- decided-by [daemon-status-skips](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-skips.md)

## History

- `01M686RDMZR9FBZXN7DZWRWQW9` Plan milestone 7: drive Rhiza from Rex and see what it does, picked by the orchestrator under the owner's standing instruction
