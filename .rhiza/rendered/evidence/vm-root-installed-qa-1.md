# vm-root-installed-qa-1

evidence · `01M6TGNN3FVJDS6VY6DT7SGVVH` · actor [claude-code-1237c200](../actor/claude-code-1237c200.md)

## result

PASS

## statement

QA claude-code-1237c200 on VM rhiza, 2026-11-07, with a binary built from the branch head (RHIZA_ROOT/CHANGE/RUN unset): `rhiza project ls` in ~/code prints nebari (initialised) and rhiza (initialised); `systemctl --user is-active rhiza-root.service` prints active and `systemctl --user status rhiza-rhiza.service` says the unit could not be found (only rhiza-root.service is in ~/.config/systemd/user); `rhiza health` in ~/code/rhiza exits 0 (daemon running as rhiza:3484327, warn only for being behind origin/main). ~/code/rhiza.toml sets self_update_project = "rhiza" and the daemon keys (maxRuns 5, permissionMode, selfUpdate, usage limits, web host); the rhiza checkout's git config keeps only rhiza.autonomy.plan, and `rhiza daemon status` reports those settings as coming from rhiza.toml.

## Links

### Out

- verifies [vm-root-installed](../capability/rhiza/orchestration/daemon/vm-root-installed.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
