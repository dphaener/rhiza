# question-99b69f45-1

question · `01M6T26ZFS8SVCP4E0JZ00RYPX`

## session

99b69f45-64e6-43c7-ad0a-01a1b4201735

## statement

Which intent should change set 01M6T1ZR (Fix the usage governor's limits) serve? milestone-29 is the only open milestone; its done point m29-one-daemon names 'one usage governor' for the root daemon. / The intent says one 10:45Z hit taught both the session and weekly limits. The code already teaches only the kind KindOf gives (reset ≤5h = session); the live hits file shows that weekly limit came from a leaked fake test hit with a 4-day reset, which change set 01M6T20E fixes. What does 01M6T1ZR do about that point? / How should a learned usage limit move up when a window's use passes it with no limit hit in that window (e.g. after the owner's upgrade to Max 20x)? / usageHoldOf adds the mean run tokens for each live run, though those runs' tokens so far are already in the window's use. How should a live run count in the projection?

## Links

### In

- answers [usage-governor-limits-fix-scope](../decision/usage-governor-limits-fix-scope.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
