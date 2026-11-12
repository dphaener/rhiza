# self-update-settings-removed

capability · `01M7642MPQYFWZPE8HSGANKMTJ` · parent [daemon-settings](../daemon-settings.md)

## acceptance

internal/rootcfg and daemon-settings tests show the old key and self_update_project reported with the autoUpdate hint.

## statement

rhiza.daemon.selfUpdate and rhiza.toml's self_update_project are removed: a root that still sets either keeps loading, selfUpdate = true is read as rhiza.daemon.autoUpdate = main (unless autoUpdate is set), the daemon logs a message naming rhiza.daemon.autoUpdate, and the next write of rhiza.toml drops the old keys; daemon set of the old key is refused with that message.

## Links

### Out

- decided-by [d-no-self-update-project](../../../../../decision/d-no-self-update-project.md)
- decided-by [d-auto-update-documented](../../../../../decision/d-auto-update-documented.md)
- decided-by [d-auto-update-name](../../../../../decision/d-auto-update-name.md)
- decided-by [d-auto-update-channels](../../../../../decision/d-auto-update-channels.md)

### In

- verifies [self-update-settings-removed-verified](../../../../../evidence/self-update-settings-removed-verified.md)
- verifies [self-update-settings-removed-qa-1](../../../../../evidence/self-update-settings-removed-qa-1.md)
- verifies [self-update-settings-removed-verified-2](../../../../../evidence/self-update-settings-removed-verified-2.md)
- verifies [self-update-settings-removed-qa-2](../../../../../evidence/self-update-settings-removed-qa-2.md)
- verifies [pr-473-auto-update-setting](../../../../../evidence/pr-473-auto-update-setting.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
