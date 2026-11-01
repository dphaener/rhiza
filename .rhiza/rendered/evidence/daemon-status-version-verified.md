# daemon-status-version-verified

evidence · `01M6APECTP7VXNBM2RARWWPAYS` · actor [claude-code-28391502](../actor/claude-code-28391502.md)

## result

PASS

## statement

At 39f8fc67, go test ./cmd/rhiza -run 'TestScripts/(daemon-status|describe-golden)$' -count=1 -v passes; the daemon-status log prints 'Version: dev' (the test build's versionString(), via ctx.Version) and the script asserts ^Version: .+$. Regenerating describe-golden with -update leaves no diff, so the golden already carries the new version field. Diff read: State gains Version recorded from ctx.Version when runDaemon starts, DaemonStatusResponse gains version (json omitempty) and a Version: line. make fmt-check, vet and test pass; self-check flags only unrendered actor files from this run's own sessions, not the leaf.

## Links

### Out

- verifies [daemon-status-version](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-version.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
