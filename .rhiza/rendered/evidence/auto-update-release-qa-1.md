# auto-update-release-qa-1

evidence · `01M77ZY705PGRA8Q3TRXXY81BN` · actor [claude-code-a0c49cdb](../actor/claude-code-a0c49cdb.md)

## result

PASS

## statement

QA at da0068fa7, hands-on with the built binary against the real published release feed (v0.1.0 release, v0.1.1-main.12 incremental): in a scratch Rhiza root with rhiza.daemon.autoUpdate=release, rhiza daemon --max-runs 0 built as v0.0.9 logged 'auto-update release: installed v0.1.0 (was v0.0.9)', re-executed, and rhiza version then reported v0.1.0. Built as v0.1.0, v0.1.0-5-gabc1234 and v0.1.1-main.5+aaa1111, each daemon ran 25s, installed nothing and left the executable byte-identical: the release channel ignores the newer incremental build and never downgrades an incremental build to the release. go test ./internal/api -run TestDaemonAutoUpdateRelease passes.

## Links

### Out

- verifies [auto-update-release](../capability/rhiza/orchestration/daemon/auto-update-release.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
