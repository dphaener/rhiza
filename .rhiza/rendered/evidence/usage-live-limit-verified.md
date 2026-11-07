# usage-live-limit-verified

evidence · `01M6TDYE4TR1R67MJPZJ0P9ZCE` · actor [claude-code-e55c799b](../actor/claude-code-e55c799b.md)

## result

PASS

## statement

Verified at 2d58ab91 by reading the diff and running scoped tests. go test -trimpath ./internal/usage/... passes: TestLiveLimit (49% of 4900 tokens gives 10000, i.e. tokens/0.49; none under 1%, with 0 or negative tokens), TestLiveWindowAndLiveness (Start = resets_at minus kind length; not live at or after resets_at). From a checkout outside ~/code, go test -trimpath ./internal/api -run 'TestUsageWindowsHitKinds|TestUsageLiveLimit|TestHealthUsage|TestUsageHold' passes. TestUsageLiveLimit uses a fake clock and fake adapter: the live 2000 beats setting 1000 and learned 900 (no hold at 1100 used); the window comes from the snapshot, with used counted to now; a lower live 1111 holds and a raised one lifts it; an expired snapshot falls back to setting 1000, then learned 900; under 1% gives no limit; a recorded limit hit still holds through governorHold. The hold is recomputed each poll (prepare keeps no hold state), so setting 1000 holding without a live limit and not holding with live 2000 shows a raised live limit lifts a hold the setting caused. TestHealthUsageLiveLimit with weekly=777 reports 'weekly 0 of 2000' and never 777. usageLimit is live, then setting, then learned, and both usageHoldOf and healthUsage call it. TestLimitHits/TestLearnedLimit* pass. Note: TestUsageWindowsHitKinds fails when run inside ~/code (the rhiza root redirects limitsDir), but the parent commit fails the same way there, so this commit did not cause it.

## Links

### Out

- verifies [usage-live-limit](../capability/rhiza/orchestration/usage-governor/usage-live-limit.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
