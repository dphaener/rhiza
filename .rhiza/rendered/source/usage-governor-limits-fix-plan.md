# usage-governor-limits-fix-plan

source · `01M6T2A9J57GY31J85KJZ2WJ3W`

## quotes

```json
[]
```

## raw_sha256

d8a7f66cf5b7ad8e6f39d68170328313f1d7a77d30c4d9453726f546b44f0aec

## session

99b69f45-64e6-43c7-ad0a-01a1b4201735

## summary

Autonomous plan run for change set 01M6T1ZR, no owner present. The planner read internal/api/usage_learn.go, usage_hold.go, internal/usage and the live daemon state: both learned limits (session 2353388, weekly 21870750) came from usage-fake test hits at 2026-11-04T10:45Z leaked into the root's usage.limits.jsonl, one of them with a 4-day reset, so the code's one-kind-per-hit rule held; the leak is change set 01M6T20E's. The planner asked four questions (which intent to serve, what to do about the both-kinds point, how learned limits relearn, how live runs count); the arbiter answered with Decision usage-governor-limits-fix-scope: serve m29-one-daemon; only a regression test for the both-kinds point; learned limits ratchet up to a window's use when it passes them with no hit in the window, settings never change; each live run adds max(0, run mean - its tokens so far) and the new pick one full mean. Three leaves follow. Nothing is left open.

## title

Usage governor limits fix, autonomous plan

## Links

### In

- derived-from [usage-governor-limits-fix-scope](../decision/usage-governor-limits-fix-scope.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
