# auto-update-release-verified-2

evidence · `01M780C9V7EKP59RGQX9W6HG1A` · actor [claude-code-8656c52d](../actor/claude-code-8656c52d.md)

## result

PASS

## statement

Verifier, attempt 2 at 26057e13f: go test -trimpath -count=1 ./internal/api -run TestDaemonAutoUpdate passes TestDaemonAutoUpdateRelease, which uses a fake release feed (httptest release.json plus tar.gz). With autoUpdate=release, v0.1.0 -> v0.2.0 installs the new executable; v0.2.0 at v0.2.0 and the incremental build v0.2.0-5-gabc1234 are left as they are; no git or make runs; latest.json is never fetched. ./internal/daemon/..., ./internal/release/... and internal/doclinks TestDocsDescribeAutoUpdate pass. The fix commit amends daemon-auto-update and auto-update-documented to name off, main, latest and release, matching daemon_settings.go and docs/runs-and-daemon.md:109 and docs/reference/config.md:50 (the major model findings f1 and f2). make self-check fails only on stale .rhiza/rendered files, which are re-rendered when the change set lands. The release test sits in internal/api, the package that wires the daemon's auto-update, not internal/daemon, which is minor as review already noted.

## Links

### Out

- verifies [auto-update-release](../capability/rhiza/orchestration/daemon/auto-update-release.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
