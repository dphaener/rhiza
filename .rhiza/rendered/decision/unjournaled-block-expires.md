# unjournaled-block-expires

decision · `01M6A0MHC4KAS2SQWFN7FFZR7K`

## statement

An unjournaled change set is held for a fixed 15 minutes (daemonUnjournaledRetry) and then picked again, re-entering the one-minute pending grace; a change set that has landed or been withdrawn leaves the unjournaled map at once; the pending grace that prevents double starts is unchanged

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m10-orchestrator-picks](../source/m10-orchestrator-picks.md)

### In

- decided-by [daemon-unjournaled-retry](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-unjournaled-retry.md)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
