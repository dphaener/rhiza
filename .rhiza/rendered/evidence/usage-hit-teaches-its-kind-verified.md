# usage-hit-teaches-its-kind-verified

evidence · `01M6T3B520J8HCKDAX5NST1PC3` · actor [claude-code-27738509](../actor/claude-code-27738509.md)

## result

PASS

## statement

At commit 4c1ee4b4 ran go test -trimpath -count=1 ./internal/api -run 'TestLimitHit' -v: TestLimitHitsBecomeLearnedLimits and TestLimitHitTeachesOnlyItsKind (subtests session_hit_over_weekly, weekly_hit_over_session) PASS. Read the diff: the new regression test seeds a learned limit of the other kind (77 tokens, hit a day earlier), records one hit with reset 3h out (session) or 4 days out (weekly) via recordLimitHit, runs learnLimits, and asserts the other kind keeps its tokens and HitAt while the hit kind is learned at the hit (350 session / 1349 weekly). Hits go to a temp Wd so none leak into a root.

## Links

### Out

- verifies [usage-hit-teaches-its-kind](../capability/rhiza/orchestration/usage-governor/usage-hit-teaches-its-kind.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
