# usage-live-run-remaining-qa-1

evidence · `01M6T65XW5QW7GW1JWVK4CRERV` · actor [claude-code-905563a6](../actor/claude-code-905563a6.md)

## result

PASS

## statement

Hands-on with a binary built at dcb78b71, scratch root, fake transcript use 5000, three finished run journals of 1000 tokens and one live run journal (process on another host), rhiza.daemon.usage.session set in git config, default 20% reserve, then rhiza daemon --once: live run at 600 with limit 7000 held at projected 6400 (5000+400+1000); live run at 1500 with limit 7000 held at projected 6000 (adds 0); live run at 1500 with limit 8000 (line 6400) did not hold, where the old count of a mean per live run (7000) would have.

## Links

### Out

- verifies [usage-live-run-remaining](../capability/rhiza/orchestration/usage-governor/usage-live-run-remaining.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
