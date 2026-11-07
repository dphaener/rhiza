# usage-live-snapshot-read-qa-1

evidence · `01M6TE9H0BGMC5YKR6WD21CKMH` · actor [claude-code-fc65330f](../actor/claude-code-fc65330f.md)

## result

PASS

## statement

QA at f6e4af6e with the built binary: the recorded usage.live.json was read back by rhiza health and rhiza daemon --once (both report the snapshot-derived limits); replacing the file with garbage made health fall back to the settings (session 0 of 555; weekly 0 of 777) with exit status unchanged and no error, so a corrupt file gives no snapshots and stops nothing.

## Links

### Out

- verifies [usage-live-snapshot-read](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-read.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
