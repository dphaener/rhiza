# usage-hold-01m6t1-qa-1

evidence · `01M6T65XYGNYV0K2HF26JEDHAF` · actor [claude-code-905563a6](../actor/claude-code-905563a6.md)

## result

PASS

## statement

Hands-on with a binary built at dcb78b71 in a scratch root: rhiza daemon --once wrote usage_hold with window, used, projected (window use + remaining live-run tokens + one mean), limit and reserve; with no setting a learned limit held, a git-config rhiza.daemon.usage.session was used as the limit, and rhiza health showed warn usage plus a usage-reserve inbox item while held. Scoped go test ./internal/api -run 'TestLimitHit|TestLearnedLimit|TestUsageHold' passes.

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
