# usage-live-limit-qa-1

evidence · `01M6TE9H9148A26139Q448CYJH` · actor [claude-code-fc65330f](../actor/claude-code-fc65330f.md)

## result

PASS

## statement

QA at f6e4af6e with the built binary under a scratch Rhiza root, one registered project, a fake transcript of 4900 counted tokens (cache reads excluded) 30 minutes old, and rhiza daemon set rhiza.daemon.usage.session 555 / weekly 777. With live snapshots 10% (5h) and 49% (7d) rhiza health showed 'session 0 of 49000; weekly 0 of 10000' (4900/0.10, 4900/0.49), not the settings. Live weekly 99% then rhiza daemon --once --permission-mode acceptEdits: the daemon held with a usage-reserve inbox item 'weekly window: 4900 tokens used, 4900 projected, over the limit of 4949 less a 20% reserve', and daemon.state.json's windows took reset = resets_at and start = resets_at minus 5h / 7d. Recording 49% and polling again cleared the hold (health 'session 4900 of 49000; weekly 4900 of 10000'). A session snapshot at 0.5% gave no live session limit (health fell back to 555) while weekly stayed live. Setting every resets_at in the past and polling fell back to the settings: hold on session at limit 555.

## Links

### Out

- verifies [usage-live-limit](../capability/rhiza/orchestration/usage-governor/usage-live-limit.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
