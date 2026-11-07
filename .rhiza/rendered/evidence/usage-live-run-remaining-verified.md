# usage-live-run-remaining-verified

evidence · `01M6T3SQPEWAZFESTVDM34VKKD` · actor [claude-code-26e102fb](../actor/claude-code-26e102fb.md)

## result

PASS

## statement

At fc4648cb ran go test -trimpath ./internal/api -run 'TestUsageHold' -count=1: TestUsageHoldProjection, TestUsageHoldStopsLaunches, TestUsageHoldStopsEveryLaunch and the new TestUsageHoldLiveRunRemaining pass. The new test uses fake journals (three finished runs of 1000, window use 100): no live run projects 1100, a live run at 600 projects 1500 (adds 400), one at 1500 projects 1100 (adds 0), the next pick adds one full mean, and at limit 1500 a live run past the mean no longer holds (the old full-mean-per-live-run count gave 2100). Read the diff: usageHoldOf now adds max(mean - usageRunTokens(st), 0) per live run plus one mean, usageRunTokens counts session and helper input+output+cache creation the same way usageRunMean does, and the result is the Projected field that usageHoldDetail prints. go test ./cmd/rhiza -run TestScripts/usage-hold also passes.

## Links

### Out

- verifies [usage-live-run-remaining](../capability/rhiza/orchestration/usage-governor/usage-live-run-remaining.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
