# auto-update-release

capability · `01M7642SFPZY7YR19WHQ5V2Y5M` · parent [daemon](../daemon.md)

## acceptance

internal/daemon auto-update tests with a fake release feed show release installing a newer release and ignoring incremental builds.

## statement

rhiza.daemon.autoUpdate=release installs the newest official release and restarts the daemon on it.

## Links

### Out

- decided-by [d-auto-update-latest](../../../../decision/d-auto-update-latest.md)
- decided-by [d-auto-update-channels](../../../../decision/d-auto-update-channels.md)
- decided-by [d-channels-planned-now](../../../../decision/d-channels-planned-now.md)

### In

- blocks [release-channels-milestone](../../../../question/release-channels-milestone.md)
- verifies [auto-update-release-verified](../../../../evidence/auto-update-release-verified.md)
- verifies [auto-update-release-qa-1](../../../../evidence/auto-update-release-qa-1.md)
- verifies [auto-update-release-verified-2](../../../../evidence/auto-update-release-verified-2.md)
- verifies [pr-489-auto-update-release-channels](../../../../evidence/pr-489-auto-update-release-channels.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
