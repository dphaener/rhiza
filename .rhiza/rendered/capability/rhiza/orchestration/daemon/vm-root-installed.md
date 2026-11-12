# vm-root-installed

capability · `01M6RN871GNRWZ6FJQAAHWT4PY` · parent [daemon](../daemon.md)

## acceptance

On the VM: rhiza project ls shows rhiza and nebari with nebari's initialised state as it is, systemctl --user shows rhiza-root.service active and no rhiza-rhiza.service, and rhiza health exits 0; the steps and their output are recorded as Evidence.

## statement

On this VM ~/code is the root with rhiza and nebari registered, rhiza-root.service runs the root daemon with rhiza.daemon.autoUpdate = main (carried over from the removed self_update_project and rhiza.daemon.selfUpdate = true), rhiza-rhiza.service is uninstalled, and the rhiza project's machine settings moved from its git config into rhiza.toml.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- decided-by [m29-registered-projects](../../../../decision/m29-registered-projects.md)

### In

- verifies [vm-root-installed-evidence](../../../../evidence/vm-root-installed-evidence.md)
- verifies [vm-root-installed-verified](../../../../evidence/vm-root-installed-verified.md)
- verifies [vm-root-installed-qa-1](../../../../evidence/vm-root-installed-qa-1.md)
- verifies [pr-351-this-vm-on-one-root](../../../../evidence/pr-351-this-vm-on-one-root.md)
- verifies [pr-473-auto-update-setting](../../../../evidence/pr-473-auto-update-setting.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
