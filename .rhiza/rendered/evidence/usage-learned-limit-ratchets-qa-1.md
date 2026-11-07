# usage-learned-limit-ratchets-qa-1

evidence · `01M6T65XT5J9JJZA883FZKRDS0` · actor [claude-code-905563a6](../actor/claude-code-905563a6.md)

## result

PASS

## statement

Hands-on with a binary built at dcb78b71, scratch root and fake Claude transcripts: learned session 1000 and weekly 2000 with no hit and window use 5000; one rhiza daemon --once saved both learned limits as 5000 in daemon.state.json, each keeping its hit_at. With a session hit inside the window (use at hit 5000, window use 7000) the session limit stayed 5000, not raised past the hit's use, while weekly rose to 7000.

## Links

### Out

- verifies [usage-learned-limit-ratchets](../capability/rhiza/orchestration/usage-governor/usage-learned-limit-ratchets.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
