# pr-349-fix-the-usage-governor-s

evidence · `01M6TGADSY2RFV5EJC92F7EAKA` · actor [dphaener](../actor/dphaener.md)

## head

cd41e7f98036712c3c1f330f7fa67e544c46df80

## pr

https://github.com/dphaener/rhiza/pull/349

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #349 (https://github.com/dphaener/rhiza/pull/349) passed CI at cd41e7f

## status

verified

## Links

### Out

- verifies [usage-live-snapshot-record](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-record.md)
- verifies [usage-live-snapshot-read](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-read.md)
- verifies [usage-live-limit](../capability/rhiza/orchestration/usage-governor/usage-live-limit.md)
- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)
- verifies [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
