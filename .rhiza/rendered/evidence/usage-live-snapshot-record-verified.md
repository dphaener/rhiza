# usage-live-snapshot-record-verified

evidence · `01M6TCC26KD6NMKF0CYCKYNK6T` · actor [claude-code-21725afd](../actor/claude-code-21725afd.md)

## result

PASS

## statement

Verifier at 8417488d: go test -trimpath ./internal/adapter/claudecode/... passes, including TestSettingsStatusLine (Settings JSON carries a statusLine command, shell-quoted, beside every hook event; none when no command is given), TestLiveUsage and TestUserStatusLine. go test -trimpath ./cmd/rhiza -run 'TestScripts/(usage-live-record|describe-golden)$' -count=1 passes. usage-live-record.txtar pipes five_hour and seven_day rate_limits into rhiza session status-line and finds session and weekly snapshots in the root's state folder (.rhiza-root/usage.live.json) and, outside a root, in .git/rhiza/usage.live.json. A second input replaces the session snapshot and leaves weekly as it was. Input without rate_limits, or malformed JSON, exits 0 and leaves the file byte-identical (cmp). With no user statusLine nothing is printed; a fake user settings statusLine 'cat' echoes stdin; a failing 'exit 3' command still exits 0. Diff read: limitsDir is the same folder the usage-limit hits use, RecordLive keys snapshots by kind and writes atomically, and session run passes StatusLineCommand [self session status-line].

## Links

### Out

- verifies [usage-live-snapshot-record](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-record.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
