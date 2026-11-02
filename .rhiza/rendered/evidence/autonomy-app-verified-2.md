# autonomy-app-verified-2

evidence · `01M6DNNT8WJ5NVEGAXSXZPPPFF` · actor [claude-code-0fed2cfc](../actor/claude-code-0fed2cfc.md)

## result

PASS

## statement

At 1a5e3269 (fix for QA f1): go test ./internal/tui/autonomyapp -count=1 -v passes all 10 tests, including plan owner->auto via daemon.set -- rhiza.autonomy.plan auto, paused toggle, maxRuns edit, start via rex window new -s S1 'rhiza daemon' -c /repo --keep-open --json -- /bin/rhiza daemon --autonomous with notice 'Starting' then 'The daemon is running' once status shows running, and the new TestStartReportsDaemonThatNeverRuns which reports 'did not start: its error is in the Rex tab' only after StartChecks reloads. go test ./cmd/rhiza -run 'TestScripts/(tui-autonomy-app|rex-install)$' passes (--once prints settings JSON; rhiza_app_autonomy on key y, no key moved). rex window new --help confirms --keep-open exists. make fmt-check vet test exit 0; make self-check fails only on stale .rhiza/rendered files from this run's model writes, no code issue. Diff read: start no longer claims success before daemon.status shows running.

## Links

### Out

- verifies [autonomy-app](../capability/rhiza/orchestration/daemon/daemon-settings/autonomy-app.md)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
