# settings-one-git-read-verified

evidence · `01M7G1PYV0N7S8MFY83TJZXAS4` · actor [claude-code-31005f93](../actor/claude-code-31005f93.md)

## result

PASS

## statement

Verified at dfca701ac. go test -trimpath -p 2 -count=1 ./internal/api -run '^TestSettingsOfOneGitRead$': exit 0 (subtests real_git, bad_value, no_keys, autonomy pass; the test counts Git calls, covers mixed-case rhiza.ship.prePushCheck, rhiza.daemon.maxRuns set twice with the last value winning, the three-part rhiza.project.watchdog.enabled with a case-different subsection kept apart, and a bad value whose usage error names the key and matches the per-key git config --get result, against fakes and a hermetic real git repo). go test -trimpath -p 2 ./internal/api -run 'Setting|WebToken|Autonomy|DaemonAutoUpdate|Health': exit 0. go test -trimpath -p 2 -timeout 8m ./internal/api: exit 0 (cache hit at the same code). go test -trimpath -p 2 -count=1 ./cmd/rhiza -run 'TestScripts/^(daemon-settings|daemon-settings-scope|daemon-load-settings|usage-settings|ship-deploy-settings|autonomy-accept-rule|stage-autonomy)$': exit 0. Measure: strace -f -qq -e trace=execve on a go test -c build of ./cmd/rhiza with -test.run 'TestScripts/^run-ship-conflict-fix-more$': PASS, 224 git config --get execs (under 400; was 1,650), 1,733 git execs in all. rhiza ship generate (scratch build): exit 0, 0 files changed; rhiza render: exit 0. Diff reviewed: settingsOf and autonomyStage/runAutonomy each read one git config -z --get-regexp ^rhiza\. snapshot per call, no cross-call cache; only internal/api source and tests plus the builder's .rhiza log and actor record changed; no TODOs or stubs.

## Links

### Out

- verifies [settings-one-git-read](../capability/settings-one-git-read.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
