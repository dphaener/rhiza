# auto-update-latest-qa-1

evidence · `01M77ZY6YXPDG1JYFGG4M8N392` · actor [claude-code-a0c49cdb](../actor/claude-code-a0c49cdb.md)

## result

PASS

## statement

QA at da0068fa7, hands-on with the built binary against the real published feeds: built rhiza stamped v0.1.1-main.5+aaa1111 into a scratch Rhiza root (isolated HOME), daemon set rhiza.daemon.autoUpdate latest (accepted; 'bogus' rejected with 'must be off or main or latest or release'), ran rhiza daemon --max-runs 0. It logged 'auto-update latest: installed v0.1.1-main.12+1f9c3a7 (was v0.1.1-main.5+aaa1111)' then 'binary changed, re-executing'; the executable's sha256 changed and rhiza version reported v0.1.1-main.12+1f9c3a7. Rebuilt stamped v0.1.1-main.12+1f9c3a7: the daemon ran 25s, logged no install, the executable was unchanged. (The re-exec'd published build predates this change and rejects the latest value, as expected for a pre-feature build.) Scoped tests: go test ./internal/daemon -run ChannelUpdate and ./internal/release pass; ./internal/api -run TestDaemonAutoUpdate tests all PASS, the package's root-state guard intermittently flags usage.live.json, which the VM's live root daemon rewrites every few seconds (the pre-existing TestDaemonAutoUpdateSetting alone trips it too).

## Links

### Out

- verifies [auto-update-latest](../capability/rhiza/orchestration/daemon/auto-update-latest.md)

## History

- `01M763ZDJBN7QKR3AAST5FTR6A` Auto-update release channels: rhiza.daemon.autoUpdate gains latest (the newest published incremental build between releases) and release (the newest official release), following the builds the release channels milestone publishes.
