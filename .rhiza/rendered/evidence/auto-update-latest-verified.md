# auto-update-latest-verified

evidence · `01M77Z3PAX26RC2QHSGB1MYNV1` · actor [claude-code-105e003c](../actor/claude-code-105e003c.md)

## result

PASS

## statement

Verifier at 38a239c92: go test -trimpath -count=1 ./internal/daemon -run ChannelUpdate passes, where a fake httptest build feed serves a latest manifest and tar.gz: TestChannelUpdateInstallsNewerBuild shows v0.1.1-main.5 replaced by the newer v0.1.1-main.7 build (exe contents and image change, so the re-exec fires, no reinstall on the next check); TestChannelUpdateLeavesUpToDateInstallAlone shows an equal or older feed version downloads nothing and leaves the exe untouched; bad checksum and unreachable feed install nothing. internal/api TestDaemonAutoUpdateLatest passes (daemon with rhiza.daemon.autoUpdate=latest installs the newer feed build over its executable, leaves an up-to-date one alone, runs no git or make), as does TestDaemonAutoUpdateSetting. internal/release, go vet and gofmt clean; cmd/rhiza testscripts daemon-settings*, describe-golden, version* pass. Read the diff: versionNewer moved to release.Newer unchanged, setting accepts latest, Download checks the manifest sha256, install writes beside the exe and renames.

## Links

### Out

- verifies [auto-update-latest](../capability/rhiza/orchestration/daemon/auto-update-latest.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
