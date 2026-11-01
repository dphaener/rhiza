# daemon-reexec-new-binary

capability · `01M6A0NRGNH76FNWH3564VWA1X` · parent [daemon](../daemon.md)

## acceptance

A Go test with an injected stat and exec: a changed executable triggers exec once, after the poll; an unchanged one never does; a pending start delays the exec to a later poll; the daemon lock is taken again by the new image

## statement

The daemon re-executes itself with the same arguments when its executable's size or modification time has changed since it started, checked at the end of a poll with no start pending, so a rebuilt rhiza takes over without a hand restart and live runs carry on

## Links

### Out

- decided-by [daemon-reexec-when-idle](../../../../decision/daemon-reexec-when-idle.md)

### In

- verifies [daemon-reexec-new-binary-verified](../../../../evidence/daemon-reexec-new-binary-verified.md)
- verifies [daemon-reexec-new-binary-qa-1](../../../../evidence/daemon-reexec-new-binary-qa-1.md)
- verifies [daemon-reexec-new-binary-verified-2](../../../../evidence/daemon-reexec-new-binary-verified-2.md)
- verifies [pr-156-daemon-picks-up-rebuilt-binary](../../../../evidence/pr-156-daemon-picks-up-rebuilt-binary.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
