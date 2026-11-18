# daemon-tmp-hold

capability · `01M7QGA6K1VD821JQCCKK11KB5` · parent [daemon-picks](../daemon-picks.md)

## acceptance

go test -trimpath ./internal/daemon/... and go test -trimpath ./internal/api -run 'Tmp|Disk|Setting|QueueWait'; a new testscript cmd/rhiza/testdata/script/daemon-tmp-low.txtar (go test -trimpath ./cmd/rhiza -run TestScripts/daemon-tmp-low) where RHIZA_TEST_TMP_FREE_BYTES (honoured only alongside RHIZA_TEST_NOW, like RHIZA_TEST_FREE_BYTES) under the minimum makes a daemon poll skip the ready change set as tmp-low naming the free space and minimum, the daemon state carry the tmp-low hold and rhiza inbox show a tmp-low item, and at or over the minimum the change set starts; (failure case) rhiza daemon set rhiza.daemon.minTmpFree banana is refused as not a size, and rhiza.daemon.minTmpFree 0 never holds; describe-golden regenerated (TestScripts/describe-golden)

## files

internal/daemon/disk.go, internal/daemon/pick.go, internal/daemon/state.go, internal/daemon/queuewait.go, internal/api/daemon.go, internal/api/daemon_settings.go, internal/api/inbox.go, internal/api/web_needs_you.go, internal/api/web_events.go, cmd/rhiza/testdata/script/daemon-tmp-low.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change the disk-low check on the repo's filesystem or rhiza.daemon.minFree; read the host's real temp free space in any test; fix the /tmp projection leak (bug 01M7M3JZFW)

## statement

The daemon's poll holds new starts and resumes while free space on its temp directory's filesystem is under rhiza.daemon.minTmpFree (default 2G, 0 off): skip reason tmp-low with the free space and the minimum in its detail, the hold saved in the daemon state, a queue-wait journal line and a tmp-low inbox item

## Links

### Out

- decided-by [m40-tmp-hold-shape](../../../../../decision/m40-tmp-hold-shape.md)

### In

- depends-on [spawn-waits-every-hold](../../../workflows/workflow-runs/spawn-waits-every-hold.md)
- verifies [daemon-tmp-hold-verified](../../../../../evidence/daemon-tmp-hold-verified.md)
- verifies [daemon-tmp-hold-qa-1](../../../../../evidence/daemon-tmp-hold-qa-1.md)
- verifies [daemon-tmp-hold-qa-2](../../../../../evidence/daemon-tmp-hold-qa-2.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
