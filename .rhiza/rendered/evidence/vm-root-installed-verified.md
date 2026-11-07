# vm-root-installed-verified

evidence · `01M6SNXF26HRA2ZBSEHKCZ32CA` · actor [claude-code-21790d75](../actor/claude-code-21790d75.md)

## result

PASS

## statement

Verifier re-ran the acceptance on VM rhiza at b8662242 (a model-only commit adding the builder's Evidence). In ~/code `rhiza project ls` exits 0 and lists nebari (initialised) and rhiza (initialised). `systemctl --user is-active rhiza-root.service` prints active; its unit runs `/home/rhiza/bin/rhiza daemon --autonomous` in /home/rhiza/code; `systemctl --user status rhiza-rhiza.service` says the unit could not be found and list-unit-files shows only rhiza-root.service. ~/code/rhiza.toml sets self_update_project = "rhiza", registers both projects and holds the [daemon] settings (maxRuns, permissionMode, selfUpdate, web host); `git config --get-regexp '^rhiza\.(daemon|web)'` in ~/code/rhiza finds none. `rhiza health` in ~/code/rhiza exits 0 with the daemon running (two warns: failed requests and unclaimed commits, unrelated).

## Links

### Out

- verifies [vm-root-installed](../capability/rhiza/orchestration/daemon/vm-root-installed.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
