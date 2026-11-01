# daemon-disk-guard

capability · `01M6ATBRV8PA51H3HXP74XNV9H` · parent [daemon-picks](../daemon-picks.md)

## acceptance

daemon-picks.txtar with RHIZA_TEST_FREE_BYTES at 1 GiB: rhiza daemon --once starts and resumes nothing and daemon status --json shows reason disk-low; at 20 GiB the same fixture starts its run; --min-free 0 disables the guard; describe-golden.txtar is regenerated

## statement

While free space on the repo's filesystem is under the daemon's minimum (10 GiB unless --min-free says otherwise), the daemon starts and resumes no runs and daemon status reports each candidate skipped as disk-low with the free space

## Links

### Out

- decided-by [disk-guard-10-gib](../../../../../decision/disk-guard-10-gib.md)

### In

- verifies [daemon-disk-guard-verified](../../../../../evidence/daemon-disk-guard-verified.md)
- verifies [daemon-disk-guard-qa-1](../../../../../evidence/daemon-disk-guard-qa-1.md)
- verifies [daemon-disk-guard-verified-2](../../../../../evidence/daemon-disk-guard-verified-2.md)
- verifies [daemon-disk-guard-qa-2](../../../../../evidence/daemon-disk-guard-qa-2.md)
- verifies [pr-162-disk-guard](../../../../../evidence/pr-162-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
