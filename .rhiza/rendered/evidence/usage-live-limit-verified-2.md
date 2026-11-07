# usage-live-limit-verified-2

evidence · `01M6TF8FRT58SWC0GYTVGR8JP7` · actor [claude-code-7368bbbc](../actor/claude-code-7368bbbc.md)

## result

PASS

## statement

Verifier at 3b30a893 (fix for the quality finding on test hermeticity): go test -trimpath -count=1 ./internal/usage/... passed; go test -trimpath -count=1 ./internal/api -run 'TestUsage|TestHealth|Limit' passed, including TestUsageLiveLimit (live limit beats the setting and the learned limit, expired snapshot falls back, raised live limit lifts a hold, recorded hit still holds) and TestHealthUsageLiveLimit (health reports weekly N of the live limit). Hermeticity: in a scratch copy under a fake root (rhiza.toml plus .rhiza-root holding an active limit hit to 2100 and 99% live snapshots), the same api tests passed and nothing was written to the fake root's state; with the fix's two test files reverted to HEAD~1, TestUsageHoldProjection, TestUsageWindowsHitKinds and TestSessionRunRecordsLimitHit failed under that root, so the usageFix/limitFixture Wd change is what keeps them out of the real root. Diff read: test-only, sets ctx.Wd outside any root.

## Links

### Out

- verifies [usage-live-limit](../capability/rhiza/orchestration/usage-governor/usage-live-limit.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
