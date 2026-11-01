# daemon-reexec-when-idle

decision · `01M6A0MKM9VJGPS6EX0JB52DW3`

## statement

At the end of a poll with nothing pending, the daemon compares its executable's size and modification time with those at startup and, if they changed, re-executes itself with the same arguments; live runs are separate processes and are untouched; the daemon records its version in its state, and daemon status prints it

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m10-orchestrator-picks](../source/m10-orchestrator-picks.md)

### In

- decided-by [daemon-reexec-new-binary](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary.md)
- decided-by [daemon-status-version](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-version.md)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
