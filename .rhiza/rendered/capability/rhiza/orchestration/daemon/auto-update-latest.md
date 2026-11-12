# auto-update-latest

capability · `01M7642SDGZMX8RYVYY562PFVR` · parent [daemon](../daemon.md)

## acceptance

internal/daemon auto-update tests with a fake build feed show latest installing a newer incremental build and leaving an up-to-date install alone.

## statement

rhiza.daemon.autoUpdate=latest installs the newest published incremental build between releases and restarts the daemon on it.

## Links

### Out

- decided-by [d-auto-update-latest](../../../../decision/d-auto-update-latest.md)
- decided-by [d-auto-update-channels](../../../../decision/d-auto-update-channels.md)
- decided-by [d-channels-planned-now](../../../../decision/d-channels-planned-now.md)

### In

- blocks [release-channels-milestone](../../../../question/release-channels-milestone.md)
- verifies [auto-update-latest-verified](../../../../evidence/auto-update-latest-verified.md)
- verifies [auto-update-latest-qa-1](../../../../evidence/auto-update-latest-qa-1.md)
- verifies [pr-489-auto-update-release-channels](../../../../evidence/pr-489-auto-update-release-channels.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
