# usage-learned-limit

capability · `01M6HJ0WFSJ9JEM9JEBHND0VX9` · parent [usage-governor](../usage-governor.md)

## acceptance

A test in internal/api's daemon tests with a fake adapter: a run driver's limit hit with a reset 3 hours out is recorded and becomes the session window's learned limit in the daemon's state, one with a reset 4 days out becomes the weekly one, and one with no reset time teaches nothing

## statement

When a run driver sees its session stop at a usage limit, it records the hit's time and reset time in the git common dir's rhiza state, which the daemon reads at each poll. The hit's window kind comes from its reset time: a reset 5 hours or less after the hit is the session window, a later one the weekly window, and a hit with no reset time teaches nothing. The daemon saves that window's use at the hit as the learned limit for its kind and uses it only when that kind has no live limit (usage-live-limit) and its rhiza.daemon.usage.<kind> setting is unset; a window with none of the three never holds

## Links

### Out

- depends-on [usage-windows](usage-windows.md)
- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- depends-on [usage-hold](usage-hold.md)
- verifies [usage-learned-limit-verified](../../../../evidence/usage-learned-limit-verified.md)
- verifies [usage-learned-limit-qa-1](../../../../evidence/usage-learned-limit-qa-1.md)
- verifies [usage-learned-limit-qa-2](../../../../evidence/usage-learned-limit-qa-2.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)
- verifies [pr-349-fix-the-usage-governor-s](../../../../evidence/pr-349-fix-the-usage-governor-s.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
