# question-c9536b48-1

question · `01M6T2HVHW6XPG83DM2D1102ST`

## session

c9536b48-1275-487b-8a2d-01a1b426e138

## statement

Where should the usage governor read the account's live usage figures from? / How should a window's limit be derived from a live snapshot? / When is a live snapshot used, and what is the fallback order? / Should this change set depend on draft 01M6T1ZR (the other usage-governor fix: a hit teaches only its kind, learned limits ratchet, live runs count only remaining tokens), which edits the same usage_hold.go/usage_learn.go code?

## Links

### In

- answers [live-usage-source](../decision/live-usage-source.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
