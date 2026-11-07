# vm-root-installed-evidence

evidence · `01M6SNSWX0XZXGAJ2303NVGFW3` · actor [claude-code-cccb9bf0](../actor/claude-code-cccb9bf0.md)

## result

PASS

## statement

On the VM: `rhiza project ls` in ~/code lists nebari (initialised) and rhiza (initialised). `systemctl --user is-active rhiza-root.service` prints active; rhiza-rhiza.service could not be found (uninstalled). `rhiza health` exits 0 in ~/code/rhiza and shows the daemon running. ~/code/rhiza.toml now sets self_update_project = "rhiza" and holds the daemon settings (maxRuns, permissionMode, selfUpdate, web host); those keys were unset from the rhiza checkout's git config, and health still passes with them.

## Links

### Out

- verifies [vm-root-installed](../capability/rhiza/orchestration/daemon/vm-root-installed.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
