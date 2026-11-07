# usage-live-run-remaining

capability · `01M6T2AQV83TNHDJ565NSM1TXG` · parent [usage-governor](../usage-governor.md)

## acceptance

go test -trimpath ./internal/api -run 'TestUsageHold' passes with tests over fake run journals: with a mean of 1000, a live run that has recorded 600 tokens adds 400 and one that has recorded 1500 adds 0, the next pick adds 1000, the hold detail's projected number is window use plus those, and a window that held before only because live runs were counted twice no longer holds

## statement

The hold's projection counts only what live runs are still expected to burn: projected use is the window's use so far, plus for each live run the mean of the last finished runs less that run's tokens so far (its journaled latest token readings, session and helpers, counted as the mean counts them), never below zero, plus one full mean for the next pick

## Links

### Out

- decided-by [usage-governor-limits-fix-scope](../../../../decision/usage-governor-limits-fix-scope.md)

### In

- verifies [usage-live-run-remaining-verified](../../../../evidence/usage-live-run-remaining-verified.md)
- verifies [usage-live-run-remaining-qa-1](../../../../evidence/usage-live-run-remaining-qa-1.md)
- verifies [pr-345-fix-the-usage-governor-s](../../../../evidence/pr-345-fix-the-usage-governor-s.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
