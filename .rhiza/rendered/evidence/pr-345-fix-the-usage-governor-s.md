# pr-345-fix-the-usage-governor-s

evidence · `01M6T6J4W2H8G0T9ACSJ7QZ3BJ` · actor [dphaener](../actor/dphaener.md)

## head

720884ebcb884f35aa964ba84e6e9182ee905696

## pr

https://github.com/dphaener/rhiza/pull/345

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #345 (https://github.com/dphaener/rhiza/pull/345) passed CI at 720884e

## status

verified

## Links

### Out

- verifies [usage-hit-teaches-its-kind](../capability/rhiza/orchestration/usage-governor/usage-hit-teaches-its-kind.md)
- verifies [usage-learned-limit-ratchets](../capability/rhiza/orchestration/usage-governor/usage-learned-limit-ratchets.md)
- verifies [usage-live-run-remaining](../capability/rhiza/orchestration/usage-governor/usage-live-run-remaining.md)
- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
