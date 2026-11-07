# vm-root-installed

capability · `01M6RN871GNRWZ6FJQAAHWT4PY` · parent [daemon](../daemon.md)

## acceptance

On the VM: rhiza project ls shows rhiza and nebari with nebari's initialised state as it is, systemctl --user shows rhiza-root.service active and no rhiza-rhiza.service, and rhiza health exits 0; the steps and their output are recorded as Evidence.

## statement

On this VM ~/code is the root with rhiza and nebari registered (rhiza as self_update_project), rhiza-root.service runs the root daemon, rhiza-rhiza.service is uninstalled, and the rhiza project's machine settings moved from its git config into rhiza.toml.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- decided-by [m29-registered-projects](../../../../decision/m29-registered-projects.md)

### In

- verifies [vm-root-installed-evidence](../../../../evidence/vm-root-installed-evidence.md)
- verifies [vm-root-installed-verified](../../../../evidence/vm-root-installed-verified.md)
- verifies [vm-root-installed-qa-1](../../../../evidence/vm-root-installed-qa-1.md)
- verifies [pr-351-this-vm-on-one-root](../../../../evidence/pr-351-this-vm-on-one-root.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
