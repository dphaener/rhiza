# pr-473-auto-update-setting

evidence · `01M77DG5R60EHW60E33SHR39QB` · actor [dphaener](../actor/dphaener.md)

## head

a5fd5a12f0d9701e6fbb89a770a6afc4faf2257d

## pr

https://github.com/dphaener/rhiza/pull/473

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #473 (https://github.com/dphaener/rhiza/pull/473) passed CI at a5fd5a1

## status

verified

## Links

### Out

- verifies [daemon-auto-update](../capability/rhiza/orchestration/daemon/daemon-auto-update.md)
- verifies [self-update-settings-removed](../capability/rhiza/orchestration/daemon/daemon-settings/self-update-settings-removed.md)
- verifies [auto-update-documented](../capability/rhiza/orchestration/daemon/daemon-settings/auto-update-documented.md)
- verifies [vm-root-installed](../capability/rhiza/orchestration/daemon/vm-root-installed.md)
- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
