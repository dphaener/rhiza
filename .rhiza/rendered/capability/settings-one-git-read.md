# settings-one-git-read

capability · `01M7FT14FW3CF1Z9Q4GW8Q33Q4`

## acceptance

A new internal/api test (go test -trimpath -p 2 -count=1 ./internal/api -run '^TestSettingsOfOneGitRead$') with a Ctx whose Git counts its calls and a repo setting several rhiza keys (one written with mixed case such as rhiza.ship.prePushCheck, one set twice so the last value wins, a three-part key such as rhiza.project.watchdog.enabled, one bad value) shows settingsOf making exactly one git config call and returning the same rows, sources and usage error naming the bad key as git config --get did. go test -trimpath -p 2 ./internal/api -run 'Setting|WebToken|Autonomy|DaemonAutoUpdate|Health' and go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/^(daemon-settings|daemon-settings-scope|daemon-load-settings|usage-settings|ship-deploy-settings|autonomy-accept-rule|stage-autonomy)$' pass. Measure: strace -f -e trace=execve on go test -run 'TestScripts/^run-ship-conflict-fix-more$' ./cmd/rhiza counts under 400 git config --get children (2026-11-15: 1,650). Touches internal/api/daemon_settings.go, internal/api/autonomy.go and their tests.

## statement

Reading the settings runs git config once per read instead of once per setting: settingsOf and the autonomy settings look every key up in one snapshot of git config with the same values, sources and validation errors as before.

## Links

### Out

- decided-by [settings-one-config-read-per-call](../decision/settings-one-config-read-per-call.md)

### In

- verifies [settings-one-git-read-verified](../evidence/settings-one-git-read-verified.md)
- verifies [settings-one-git-read-qa-1](../evidence/settings-one-git-read-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
