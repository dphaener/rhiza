# disk-guard-10-gib

decision · `01M6ATB21VA3XW8EKPA0SAAEKZ`

## statement

While free space on the repo's filesystem is under 10 GiB (overridable with daemon --min-free), the daemon starts and resumes no runs, reports each candidate skipped with reason disk-low naming the free space, and the inbox shows one disk-low item read from the daemon's state; tests inject free space through RHIZA_TEST_FREE_BYTES, honoured only when RHIZA_TEST_NOW is set

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m12-orchestrator-picks](../source/m12-orchestrator-picks.md)

### In

- decided-by [daemon-disk-guard](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-disk-guard.md)
- decided-by [inbox-disk-low](../capability/rhiza/orchestration/inbox/inbox-list/inbox-disk-low.md)

## History

- `01M6ATAR49P8K7TM7RCW7D20KF` Plan milestone 12: the night runs clean, picked by the orchestrator under the owner's standing instruction
