# usage-hit-teaches-its-kind-qa-1

evidence · `01M6T65XKNP98TD08XF2855SFZ` · actor [claude-code-905563a6](../actor/claude-code-905563a6.md)

## result

PASS

## statement

Hands-on with a binary built at dcb78b71, in a scratch root (one project, isolated HOME and CLAUDE_CONFIG_DIR, fake Claude transcripts with 5000 then 2000 counted tokens): daemon state learned session 1000 and weekly 2000 (hit 3 days ago), one session-kind hit recorded in .rhiza-root/usage.limits.jsonl 30 min ago with reset 150 min out. After rhiza daemon --once the session limit was taught at the hit (5000, new hit time) while the weekly learned limit kept its old hit time and was not taught from the hit (it only ratcheted to its window use 7000).

## Links

### Out

- verifies [usage-hit-teaches-its-kind](../capability/rhiza/orchestration/usage-governor/usage-hit-teaches-its-kind.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
