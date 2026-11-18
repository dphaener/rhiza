# daemon-tmp-hold-qa-1

evidence · `01M7QK1TTVYK0PF2P50JPY3MQH` · actor [claude-code-062a7236](../actor/claude-code-062a7236.md)

## result

PASS

## statement

QA with a scratch build in a hermetic repo: RHIZA_TEST_TMP_FREE_BYTES under rhiza.daemon.minTmpFree made daemon --once skip the ready change set as tmp-low naming free space and minimum, saved tmp_low in the daemon state, wrote a queue-wait line and one tmp-low inbox item; at the minimum, with 0, or with --min-tmp-free lowered it started; banana and 2Q refused as not a size

## Links

### Out

- verifies [daemon-tmp-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-tmp-hold.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
