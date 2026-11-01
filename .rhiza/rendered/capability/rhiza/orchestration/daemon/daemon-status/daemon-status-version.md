# daemon-status-version

capability · `01M6A0NSV7XVE3SXVN0V6AKTKR` · parent [daemon-status](../daemon-status.md)

## acceptance

daemon-status.txtar shows a Version: line with the test build's version; describe-golden.txtar is regenerated for the new field

## statement

rhiza daemon status prints the version of the binary the daemon runs, recorded in the daemon's state when it starts or re-executes

## Links

### Out

- decided-by [daemon-reexec-when-idle](../../../../../decision/daemon-reexec-when-idle.md)

### In

- verifies [daemon-status-version-verified](../../../../../evidence/daemon-status-version-verified.md)
- verifies [daemon-status-version-qa-1](../../../../../evidence/daemon-status-version-qa-1.md)
- verifies [pr-156-daemon-picks-up-rebuilt-binary](../../../../../evidence/pr-156-daemon-picks-up-rebuilt-binary.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
