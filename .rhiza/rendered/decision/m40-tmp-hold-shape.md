# m40-tmp-hold-shape

decision · `01M7QGA6F5QRHQQ8XHJP26P5BJ`

## statement

The free-tmp check is a daemon poll hold like disk-low: free space on the filesystem of the daemon's temp directory (os.TempDir, so TMPDIR when set) under rhiza.daemon.minTmpFree, a size defaulting to 2G with 0 turning it off, holds new starts and resumes with skip reason tmp-low naming the free space and the minimum

## wording

(assumed by the autonomous planner from the existing disk-low guard and Decision m40-host-limits-as-holds)

## Links

### Out

- derived-from [m40-spawn-holds-plan](../source/m40-spawn-holds-plan.md)

### In

- decided-by [daemon-tmp-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-tmp-hold.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
