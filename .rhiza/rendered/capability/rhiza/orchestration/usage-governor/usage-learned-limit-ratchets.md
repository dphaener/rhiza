# usage-learned-limit-ratchets

capability · `01M6T2AQSEY1KKA7VTESYMXWD8` · parent [usage-governor](../usage-governor.md)

## acceptance

go test -trimpath ./internal/api -run 'TestLearnedLimit' passes with tests using a fake clock and fake token records: a window whose use passes the learned weekly limit with no hit raises it to that use and saves it in daemon state; a window under the learned limit leaves it; a window holding a hit of that kind does not raise it past the hit's use; a later hit with lower use moves it down; and with rhiza.daemon.usage.weekly set, usageLimit still returns the setting

## statement

A learned limit only ratchets up between hits: at each poll, when a current window's use passes its kind's learned limit and no recorded hit of that kind falls inside the window, the daemon saves that use as the kind's learned limit (keeping the hit time it had), so limits learned before a plan upgrade rise with use; a later hit still replaces the learned limit with the use at that hit, up or down, and a limit setting is never changed

## Links

### Out

- decided-by [usage-governor-limits-fix-scope](../../../../decision/usage-governor-limits-fix-scope.md)

### In

- verifies [usage-learned-limit-ratchets-verified](../../../../evidence/usage-learned-limit-ratchets-verified.md)
- verifies [usage-learned-limit-ratchets-qa-1](../../../../evidence/usage-learned-limit-ratchets-qa-1.md)
- verifies [pr-345-fix-the-usage-governor-s](../../../../evidence/pr-345-fix-the-usage-governor-s.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
