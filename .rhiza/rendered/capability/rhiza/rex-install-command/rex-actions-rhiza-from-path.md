# rex-actions-rhiza-from-path

capability · `01M6JBWYEG2NS37GRDPKMSWJ7D` · parent [rex-install-command](../rex-install-command.md)

## acceptance

A test installs with ctx.Self set to an absolute path and finds that path nowhere in rhiza.lua; rhiza.lua takes BIN from settings.bin only when it is a non-empty string (type(settings.bin) == "string" and settings.bin ~= "", as launch_key is read) and otherwise uses "rhiza", so a missing, empty or non-string bin yields "rhiza"; a reinstall keeps a bin the owner set in the settings block byte for byte, and a second install is byte-identical

## statement

The app actions rhiza rex install writes run rhiza from the PATH of the host running the pane, or from an optional bin entry in rhiza.lua's settings block, never from a path fixed at install time

## Links

### Out

- decided-by [rex-actions-run-rhiza-from-path](../../../decision/rex-actions-run-rhiza-from-path.md)

### In

- verifies [rex-actions-rhiza-from-path-verified](../../../evidence/rex-actions-rhiza-from-path-verified.md)
- verifies [rex-actions-rhiza-from-path-qa-1](../../../evidence/rex-actions-rhiza-from-path-qa-1.md)
- verifies [pr-266-fix-rhiza-rex-install-bakes](../../../evidence/pr-266-fix-rhiza-rex-install-bakes.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
