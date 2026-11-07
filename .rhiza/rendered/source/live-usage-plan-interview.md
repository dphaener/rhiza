# live-usage-plan-interview

source · `01M6T2NE1PP99H43ZG609Z6W2W`

## quotes

```json
[]
```

## raw_sha256

05dc6576c05c762d0ccb3b9eb351106f3529d8db3c7d26439c4254b9a22bf6b7

## session

c9536b48-1275-487b-8a2d-01a1b426e138

## summary

Autonomous plan interview for change set 01M6T2DC (usage governor limits are static). No owner present. The planner read the usage governor (internal/usage, internal/api/usage_hold.go, usage_learn.go), the claudecode adapter's --settings JSON and the sibling draft 01M6T1ZR, and put four questions to the arbiter: where live usage comes from, how a window's limit is derived, the fallback order, and ordering against 01M6T1ZR. The arbiter answered with Decision live-usage-source: Rhiza-launched Claude Code sessions get a statusLine command that records rate_limits (used_percentage, resets_at, capture time) per window kind into the limit-hits state folder, read through an optional adapter interface, no network or credentials; a live snapshot fixes its window (reset = resets_at, start = reset minus length) and the limit is the window's counted tokens up to capture divided by used_percentage/100, none below 1% or with no tokens; a snapshot is live while resets_at is ahead; order per kind is live, setting, learned; the limit-hit hold is unchanged; this change set lands after 01M6T1ZR. The planner then wrote three leaves under usage-governor and the change set's serves link to m29-one-daemon.

## title

Live usage limits plan interview

## Links

### In

- derived-from [live-usage-source](../decision/live-usage-source.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
