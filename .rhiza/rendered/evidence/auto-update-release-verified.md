# auto-update-release-verified

evidence · `01M77Z9BCTYDBSBV7XWSXG9A04` · actor [claude-code-a09a2bbb](../actor/claude-code-a09a2bbb.md)

## result

PASS

## statement

At 24073c15a, go test -trimpath -count=1 ./internal/api -run TestDaemonAutoUpdate -v passes, including the new TestDaemonAutoUpdateRelease: with rhiza.daemon.autoUpdate=release and a fake httptest release feed, a daemon on v0.1.0 installs the v0.2.0 release over its executable (via daemon.ChannelUpdate on release.ChannelRelease, no git or make), an up-to-date v0.2.0 is left alone, an incremental build v0.2.0-5-gabc1234 is not replaced by the release, and the latest feed is never fetched (the fake fails the test if it is). internal/daemon and internal/release tests pass. The diff wires the release value in runDaemon to release.ChannelRelease, accepts release in the setting's validation, and updates docs/reference/config.md and docs/runs-and-daemon.md. Note: the test sits in internal/api (daemon auto-update tests), not internal/daemon, where the existing latest-channel tests also live. make self-check only fails on stale .rhiza/rendered actor files from this run, fixed by rhiza render at close-out.

## Links

### Out

- verifies [auto-update-release](../capability/rhiza/orchestration/daemon/auto-update-release.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
