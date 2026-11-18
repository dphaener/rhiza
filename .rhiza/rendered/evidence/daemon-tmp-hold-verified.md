# daemon-tmp-hold-verified

evidence · `01M7QHBJFRASWH6EQGJZ0CTM0S` · actor [claude-code-9121e723](../actor/claude-code-9121e723.md)

## result

PASS

## statement

At 921336c28, re-ran every Checks: command in capped scopes, each exit 0: go test -trimpath -p 2 ./internal/daemon/...; ./internal/api -run 'Tmp|Disk|Setting|QueueWait'; ./internal/api whole (180s); ./cmd/rhiza TestScripts/daemon-tmp-low (-v: PASS); TestScripts/describe-golden -update (golden unchanged, tree clean after); the 23-script daemon/settings list; TestScripts/run-ship-change; rhiza ship generate (0 changed files) and rhiza render (only the verifier's own actor page). Diff read: daemonTmpLow/tmpFreeBytes read TMPDIR (else os.TempDir) and, under RHIZA_TEST_NOW, only RHIZA_TEST_TMP_FREE_BYTES (unknown otherwise, so no test reads host tmp); the hold goes to daemonPoll1 and schedules/demand via daemonConf.hostHold as skip reason tmp-low with the detail '1.0 GiB free in the temp directory, under the minimum of 2.0 GiB', is saved as tmp_low in daemon.state.json, is in queue-wait's hold reasons, and gives a machine-wide tmp-low inbox item. rhiza.daemon.minTmpFree has a 2G default, is validated as a size, and has a --min-tmp-free flag. daemon-tmp-low.txtar covers: under the minimum the change set is skipped with the state, journal and inbox checked; 'banana' is refused as not a size; 0 never holds; at the minimum the change set starts. The disk-low/minFree logic is unchanged, and the files stay within the leaf with no TODOs or stubs

## Links

### Out

- verifies [daemon-tmp-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-tmp-hold.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
