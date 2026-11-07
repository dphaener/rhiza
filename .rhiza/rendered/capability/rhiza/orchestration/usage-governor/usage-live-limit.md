# usage-live-limit

capability · `01M6T2P2PWNPYM6T54QB5Q08K6` · parent [usage-governor](../usage-governor.md)

## acceptance

Unit tests in internal/usage derive window and limit from a snapshot (49% of a 7-day window's counted tokens gives tokens/0.49), give no limit under 1% or with no tokens, and treat a snapshot past its resets_at as not live; internal/api usage_hold tests with a fake clock and fake adapter show a live limit wins over a set rhiza.daemon.usage.weekly and a learned limit, an expired snapshot falls back to the setting and then the learned limit, a raised live limit lifts a hold the old setting caused, and a recorded limit hit still holds regardless; an internal/api health test with the same fake live snapshot and a set rhiza.daemon.usage.weekly shows the usage part reporting 'weekly N of <live limit>', not the setting

## statement

At each poll a window kind with a live snapshot (its resets_at after now) takes its window from it, Reset = resets_at and Start = resets_at minus the kind's length, and its limit from it, the window's counted tokens over [Start, capture time) divided by used_percentage/100; a snapshot under 1% used or with no tokens in its window gives no limit. Each kind's limit is the live one, else the rhiza.daemon.usage.<kind> setting, else the learned limit, and with none the kind never holds; Used is still counted fresh up to now, and the machine-wide hold from a recorded limit hit is unchanged; every caller that needs a kind's limit, the poll's hold and rhiza health's usage part alike, goes through this one live-first limit, so health reports each window's use against the same limit that holds work

## Links

### Out

- depends-on [usage-live-snapshot-read](usage-live-snapshot-read.md)
- decided-by [live-usage-source](../../../../decision/live-usage-source.md)

### In

- verifies [usage-live-limit-verified](../../../../evidence/usage-live-limit-verified.md)
- verifies [usage-live-limit-qa-1](../../../../evidence/usage-live-limit-qa-1.md)
- verifies [usage-live-limit-verified-2](../../../../evidence/usage-live-limit-verified-2.md)
- verifies [pr-349-fix-the-usage-governor-s](../../../../evidence/pr-349-fix-the-usage-governor-s.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
