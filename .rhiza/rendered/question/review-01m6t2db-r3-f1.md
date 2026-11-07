# review-01m6t2db-r3-f1

question · `01M6T3CQEQY8XWPYSVBDCEXDNW`

## statement

Plan review round 3 left finding f1 (contradiction on usage-governor-hold-rules): The plan rewrites this Decision's statement and it still says 'The daemon holds a new pick when the window's use so far, plus the mean tokens of the last 10 finished runs for each live run and for the one it would start, exceeds the limit minus the reserve'. That counts a full mean for each live run. Draft 01M6T1ZR (accepted, which this change set lands after) adds usage-live-run-remaining: each live run adds only the mean less its tokens so far, never below zero. This plan's own rewritten usage-hold says the same, citing usage-live-run-remaining. 01M6T1ZR's ops leave this Decision alone, so once both land the rewritten text contradicts usage-live-run-remaining and usage-hold. The rewrite already folds in 01M6T1ZR's ratchet; the live-run term should match too: the mean less each live run's tokens so far, never below zero, plus one full mean for the next pick. The critic confirmed it: The plan's rewritten usage-governor-hold-rules (01M6HHVZ) statement still says the hold counts 'the mean tokens of the last 10 finished runs for each live run and for the one it would start', a full mean per live run. Accepted draft 01M6T1ZR adds usage-live-run-remaining, under which each live run adds the mean less its tokens so far, never below zero, plus one full mean for the next pick. This plan's own usage-hold says the same. 01M6T1ZR's ops never touch 01M6HHVZ, so once both land the Decision contradicts usage-live-run-remaining and usage-hold. The rewrite already folds in 01M6T1ZR's ratchet, so it should fold in this rule too. How should it be settled?

## Links

### In

- answers [usage-hold-rules-live-run-remaining](../decision/usage-hold-rules-live-run-remaining.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
