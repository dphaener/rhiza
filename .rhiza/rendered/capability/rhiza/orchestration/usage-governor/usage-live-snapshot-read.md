# usage-live-snapshot-read

capability · `01M6T2P2P020GWJEF7T5JE0FTY` · parent [usage-governor](../usage-governor.md)

## acceptance

internal/adapter/claudecode tests read a fake state folder's live-usage file back as the latest snapshot per kind with used_percentage, resets_at and capture time, and an absent or corrupt file gives no snapshots and no error that stops a poll

## statement

An adapter that reports live usage implements an optional interface, like UsageReader, that returns the latest recorded snapshot per window kind; the Claude Code adapter reads the live-usage file, and an adapter without it reports none

## Links

### Out

- depends-on [usage-live-snapshot-record](usage-live-snapshot-record.md)
- decided-by [live-usage-source](../../../../decision/live-usage-source.md)

### In

- depends-on [usage-live-limit](usage-live-limit.md)
- verifies [usage-live-snapshot-read-verified](../../../../evidence/usage-live-snapshot-read-verified.md)
- verifies [usage-live-snapshot-read-qa-1](../../../../evidence/usage-live-snapshot-read-qa-1.md)
- verifies [pr-349-fix-the-usage-governor-s](../../../../evidence/pr-349-fix-the-usage-governor-s.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
