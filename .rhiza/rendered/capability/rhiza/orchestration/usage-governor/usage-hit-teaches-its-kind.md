# usage-hit-teaches-its-kind

capability · `01M6T2AQQYD3ESRA7NZAQW22KT` · parent [usage-governor](../usage-governor.md)

## acceptance

go test -trimpath ./internal/api -run 'TestLimitHit' passes with a regression test that records one session-kind hit (reset 3 hours out) over an already learned weekly limit and asserts the weekly learned limit keeps its tokens and hit time, and the reverse for a weekly-kind hit over a learned session limit

## statement

A usage-limit hit teaches only the window kind its reset time gives: a session-kind hit leaves the learned weekly limit as it was, and a weekly-kind hit leaves the learned session limit as it was

## Links

### Out

- decided-by [usage-governor-limits-fix-scope](../../../../decision/usage-governor-limits-fix-scope.md)

### In

- verifies [usage-hit-teaches-its-kind-verified](../../../../evidence/usage-hit-teaches-its-kind-verified.md)
- verifies [usage-hit-teaches-its-kind-qa-1](../../../../evidence/usage-hit-teaches-its-kind-qa-1.md)
- verifies [pr-345-fix-the-usage-governor-s](../../../../evidence/pr-345-fix-the-usage-governor-s.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
