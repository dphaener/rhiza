# daemon-auto-update

capability · `01M7642MMAZSKZD27543JQ79D8` · parent [daemon](../daemon.md)

## acceptance

internal/daemon self-update tests with fake git and make cover off, main with a Rhiza checkout, and main with none (logged, nothing updated); the daemon-settings testscript shows autoUpdate=off by default; the setting accepts off, main, latest and release and rejects any other value; the latest and release channels are checked by auto-update-latest's and auto-update-release's acceptance.

## statement

rhiza.daemon.autoUpdate takes off (the default), main, latest or release: with main the daemon finds the registered project that is Rhiza's own source checkout by its module path, fast-forwards it when origin/main moves and runs make update, as self-update does today; latest and release install the newest published build of their channel and restart the daemon on it (auto-update-latest, auto-update-release); off never updates.

## Links

### Out

- decided-by [d-no-self-update-project](../../../../decision/d-no-self-update-project.md)
- decided-by [d-auto-update-documented](../../../../decision/d-auto-update-documented.md)
- decided-by [d-auto-update-name](../../../../decision/d-auto-update-name.md)
- decided-by [d-auto-update-channels](../../../../decision/d-auto-update-channels.md)

### In

- verifies [daemon-auto-update-verified](../../../../evidence/daemon-auto-update-verified.md)
- verifies [daemon-auto-update-qa-1](../../../../evidence/daemon-auto-update-qa-1.md)
- verifies [daemon-auto-update-verified-2](../../../../evidence/daemon-auto-update-verified-2.md)
- verifies [daemon-auto-update-qa-2](../../../../evidence/daemon-auto-update-qa-2.md)
- verifies [pr-473-auto-update-setting](../../../../evidence/pr-473-auto-update-setting.md)
- verifies [pr-489-auto-update-release-channels](../../../../evidence/pr-489-auto-update-release-channels.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
