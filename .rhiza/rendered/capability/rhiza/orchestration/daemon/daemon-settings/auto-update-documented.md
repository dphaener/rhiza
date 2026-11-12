# auto-update-documented

capability · `01M7642MRZ4256B2ECYVG4D4SH` · parent [daemon-settings](../daemon-settings.md)

## acceptance

The generated settings reference golden is regenerated and shows autoUpdate as off, main, latest or release; docs/runs-and-daemon.md and docs/reference/config.md describe latest and release as working channels; a docs test (go test ./internal/doclinks -run TestDocsDescribeAutoUpdate) fails unless docs/runs-and-daemon.md names rhiza.daemon.autoUpdate with its off and main values and no longer names selfUpdate or self_update_project.

## statement

rhiza.daemon.autoUpdate and its channels are described in daemon settings, the generated settings reference and the user docs on the daemon as off, main, latest or release, each a working channel with what it updates from.

## Links

### Out

- decided-by [d-no-self-update-project](../../../../../decision/d-no-self-update-project.md)
- decided-by [d-auto-update-documented](../../../../../decision/d-auto-update-documented.md)
- decided-by [d-auto-update-name](../../../../../decision/d-auto-update-name.md)
- decided-by [d-auto-update-channels](../../../../../decision/d-auto-update-channels.md)

### In

- verifies [auto-update-documented-verified](../../../../../evidence/auto-update-documented-verified.md)
- verifies [auto-update-documented-qa-1](../../../../../evidence/auto-update-documented-qa-1.md)
- verifies [auto-update-documented-qa-2](../../../../../evidence/auto-update-documented-qa-2.md)
- verifies [pr-473-auto-update-setting](../../../../../evidence/pr-473-auto-update-setting.md)
- verifies [pr-489-auto-update-release-channels](../../../../../evidence/pr-489-auto-update-release-channels.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
