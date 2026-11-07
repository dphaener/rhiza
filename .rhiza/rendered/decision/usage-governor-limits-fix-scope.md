# usage-governor-limits-fix-scope

decision · `01M6T298K6HQ3H61GT3Y0EACCE`

## reasoning

Read the question, @usage-learned-limit, internal/api/usage_learn.go, usage_hold.go and usage.KindOf. m29-one-daemon is the only open milestone's done point that names the governor, so the change set serves it. learnLimits keys each hit by usage.KindOf and keeps the latest per kind, so one hit cannot teach both kinds; the owner's memory notes and the planning request show every recorded hit was a usage-fake test hit, already a separate change set, so duplicating that fix here would overlap it; a regression test pins the behaviour the intent asks for. For relearning, a window passing its limit without a hit proves the real limit is at least that use, so raising to observed use is safe and never invents headroom; it rises slowly because the hold stops new picks near the limit, which is acceptable since an explicit setting (as the owner now has in rhiza.toml) overrides it. For live runs, the window use already includes their tokens so far, so only the rest of a mean run is still to come; clamping at zero keeps a long run from lowering the projection.

## statement

(1) Change set 01M6T1ZR serves @m29-one-daemon (milestone-29's 'one usage governor'). (2) On a hit teaching both limits: no new code here; the code already teaches only the kind KindOf gives, and the wrong weekly limit came from a leaked fake test hit, which the fake-usage-hits fix (request 01M6T20A, run 01M6T20C) removes. 01M6T1ZR only adds a regression test: a session-kind hit leaves the learned weekly limit unchanged. (3) A learned limit is a lower bound and only ratchets up: at each poll, when a window's use passes its kind's learned limit and no hit falls in that window, the learned limit becomes that use (keeping its kind; a later hit replaces it as now, and may move it down). A limit setting is never changed; setting it stays the way to jump straight to a new plan's limit. (4) The projection counts each live run's remaining expected tokens only: projected = window use + for each live run max(0, run mean - that run's tokens so far) + one run mean for the new pick.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-99b69f45-1](../question/question-99b69f45-1.md)
- derived-from [usage-governor-limits-fix-plan](../source/usage-governor-limits-fix-plan.md)

### In

- decided-by [usage-hit-teaches-its-kind](../capability/rhiza/orchestration/usage-governor/usage-hit-teaches-its-kind.md)
- decided-by [usage-learned-limit-ratchets](../capability/rhiza/orchestration/usage-governor/usage-learned-limit-ratchets.md)
- decided-by [usage-live-run-remaining](../capability/rhiza/orchestration/usage-governor/usage-live-run-remaining.md)
- decided-by [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
