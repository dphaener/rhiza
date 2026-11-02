# daemon-writes-state-at-start

decision · `01M6CA9ZT2YKEQ41XMMKB91M1R`

## statement

A starting daemon image writes its state once before its first poll, with its own process, start time and settings and the previous image's last poll, live runs, skips and unjournaled holds, so status matches the lock holder at once

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m21-orchestrator-picks](../source/m21-orchestrator-picks.md)

### In

- decided-by [daemon-state-at-start](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary/daemon-state-at-start.md)

## History

- `01M6CA9PFRY55ZTXNHWDM95Y9A` Plan milestone 21: status tells the truth, picked by the orchestrator under the owner's standing instruction
