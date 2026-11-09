# question-375c0132-1

question · `01M6XVBG8ZBCCBNAE97KP82RC0`

## session

375c0132-d622-41bb-97c1-01a1bbb2b95d

## statement

Change set 01M6XV4 (plan runs exceed rhiza.daemon.maxRuns): the root cause found is that each project's ship picks, continuations and queued starts count only live runs, not a plan request's start that has not journaled yet (up to daemonStartGrace, 1m), while the round-robin share already assumed it holds a slot; at 13:16:15 a request's plan run started and at 13:16:36 two ship runs started before it journaled. Should plan runs share the machine-wide maxRuns, or get their own cap? / The daemon's own resume of an interrupted run (daemon.Resumable) starts without any slot check and is not counted by later picks in the same poll. Should it count against maxRuns? / A milestone request's plan run currently starts at once whether or not a slot is free (owner decision milestone-interview-starts-at-once, an interactive owner session). Keep that exception? / How should rhiza health show live runs against the cap?

## Links

### In

- answers [plan-runs-share-max-runs](../decision/plan-runs-share-max-runs.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
