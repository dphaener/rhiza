# usage-live-snapshot-record-qa-1

evidence · `01M6TE9GTDNC6M11XW0X84Y35B` · actor [claude-code-fc65330f](../actor/claude-code-fc65330f.md)

## result

PASS

## statement

QA at f6e4af6e with a binary built from the branch, in a hermetic scratch repo (own HOME, RHIZA_* unset). rhiza session status-line fed Claude Code-shaped JSON (session_id, model, rate_limits five_hour 10% and seven_day 49%): exit 0, printed the user's own statusLine output (a printf+head command from HOME/.claude/settings.json, receiving the same stdin), --json said recorded 2, and .git/rhiza/usage.live.json held one session and one weekly snapshot with used_percentage, resets_at and captured_at. After rhiza root init on the parent dir the same command wrote to <root>/.rhiza-root/usage.live.json instead; a later input with only seven_day replaced weekly and left session. rhiza session start (fake rex) then rhiza session run <id> --dry-run --json: the claude argv's --settings JSON carries statusLine {type: command, command: '<rhiza path> session status-line'} beside the 7 hook events.

## Links

### Out

- verifies [usage-live-snapshot-record](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-record.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
