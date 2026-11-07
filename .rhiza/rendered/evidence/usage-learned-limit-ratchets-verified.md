# usage-learned-limit-ratchets-verified

evidence · `01M6T5V5513KHXY43V6Z07FYRX` · actor [claude-code-c9f9418e](../actor/claude-code-c9f9418e.md)

## result

PASS

## statement

At 6e29d497, `go test -trimpath -count=1 ./internal/api -run 'TestLearnedLimit' -v` passes. TestLearnedLimitRatchets checks raiseLearned directly. TestLearnedLimitRatchetsAtPoll drives daemonProj.prepare and saveState with a fake ctx.Now and fake token records (holdFakeRecords), then reads the learned limits back with daemon.ReadState. It covers four cases. Use 1500 over a learned 1000 with no hit is saved as 1500 and keeps the old hit time. Use under a learned 2000 leaves it. A window holding a weekly hit saves the 1500 used at the hit, not the 2200 used since, with the new hit time. After the limit is raised to 1500, a later week's hit at 1200 use moves it down to 1200, and with rhiza.daemon.usage.weekly=42 set, usageLimit still returns 42. The diff calls learnLimits before ratchetLimits in prepare. ratchetLimits leaves alone any kind whose window holds a hit of that kind, and it never touches settings. gofmt and go vet on internal/api are clean. All three findings from attempt 1 are addressed. The wider run `-run 'TestLimitHit|TestUsage|TestLearn'` passes from a worktree outside ~/code. Under ~/code, TestUsageWindowsHitKinds fails: it predates this leaf and is unchanged by it, and it reads leaked usage hits from the ~/code root's state.

## Links

### Out

- verifies [usage-learned-limit-ratchets](../capability/rhiza/orchestration/usage-governor/usage-learned-limit-ratchets.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
