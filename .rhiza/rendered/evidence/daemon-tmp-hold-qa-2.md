# daemon-tmp-hold-qa-2

evidence · `01M7QKYRB7YT1JBSWS34V3XFBD` · actor [claude-code-608090ea](../actor/claude-code-608090ea.md)

## result

PASS

## statement

QA cycle 2, scratch build of 13a633dfe in a hermetic repo (RHIZA_TEST_NOW stand-ins, fake rex): RHIZA_TEST_TMP_FREE_BYTES=1GiB under the 2G default made daemon --once skip the ready change set as tmp-low ('1.0 GiB free in the temp directory, under the minimum of 2.0 GiB'), saved tmp_low in daemon.state.json, wrote a tmp-low queue-wait line and one tmp-low inbox item; 2147483647 bytes held, exactly 2G started; minTmpFree 0 started even at 1 byte and cleared the item and state; daemon set minTmpFree banana and 2Q and --min-tmp-free banana refused as not a size (exit 2)

## Links

### Out

- verifies [daemon-tmp-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-tmp-hold.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
