# plan-runs-share-max-runs

decision · `01M6XVDZHN24987TJWZ6FNEHXD`

## reasoning

The change set's intent is the owner's own report: load 23 on a 4-core VM with 5 runs at cap 3 and 4, and it states plan runs 'should count against the machine-wide cap' and a waiting request 'may reserve the next free slot (#226) but must not exceed the cap'; it offers a separate plan cap only as an alternative, and one shared cap is simpler, matches 'machine-wide' and is what bounds load, which is the harm. The question's root cause (daemon_projects.go counts pending starts within daemonStartGrace per project, while plan starts from request.go were not seen by the other picks) means the fix is to count every unjournaled start everywhere, not a new cap. A resume of an interrupted run (daemon.go:827, daemon.Resumable) is a full run using agents and CPU like any other, so leaving it uncapped and uncounted repeats the same overshoot; resuming first keeps interrupted work ahead of new work. The milestone exception was itself an arbiter assumption (wording assumed), reasoned on the interview costing 'little compute'; but a plan run still starts agents, it was one of the plan runs in the reported burst, and the owner's later explicit words 'must not exceed the cap' outrank that assumption. Giving the interview first claim on the next slot keeps the point of that decision (an owner interview is not queued behind build picks) without breaking the cap. Health already reports per-part ok/warn/bad lines (health_schedules.go); live/max with starts listed separately is what the intent asks for ('show live runs against the cap') and makes this bug visible as bad when it recurs.

## statement

Every run the daemon starts counts against the one machine-wide rhiza.daemon.maxRuns, and nothing it starts may take running past it. (1) Plan runs share maxRuns; there is no separate plan cap. Every pick (ship picks, continuations, queued starts, request plan runs, schedules) counts live runs plus every start it or another project made that has not journaled yet and is still within daemonStartGrace, across all projects of the daemon, so a plan run started at 13:16:15 holds its slot when the 13:16:36 picks run. A waiting request keeps reserving the next free slot (#226) but never exceeds the cap. (2) The daemon's resume of an interrupted run (daemon.Resumable) counts against maxRuns: it starts only when a slot is free, resumes go first among the poll's starts, and each resume counts as a pending start for later picks in the same poll; a resume that finds no slot waits for a later poll. (3) This supersedes the part of milestone-interview-starts-at-once that let a milestone interview start without a free slot: a queued milestone request's plan run takes the next free slot ahead of every other start (resumes included) and every other request, but waits while running is at the cap. The rest of that decision stands: one milestone interview live at a time, not held by a usage-reserve hold, held by pause and disk-low, waiting status and inbox item at launch. (4) rhiza health shows a runs line: live/max with pending starts shown separately (e.g. 'runs: 3/4 live, 1 starting'), summed across the daemon's projects; the part is ok at or under the cap and bad when live plus starting exceeds maxRuns, naming the runs over the cap. The daemon status (daemonapp/web) shows the same live/max count.

## wording

(assumed, would have escalated)

## Links

### Out

- relates-to [milestone-interview-starts-at-once](milestone-interview-starts-at-once.md) — supersedes its start without a free slot; the rest of that decision stands
- answers [question-375c0132-1](../question/question-375c0132-1.md)
- derived-from [plan-runs-cap-interview](../source/plan-runs-cap-interview.md)

### In

- decided-by [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)
- decided-by [daemon-resume-takes-slot](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-resume-takes-slot.md)
- decided-by [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- decided-by [request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/request-slot-kept.md)
- decided-by [health-runs-against-cap](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-runs-against-cap.md)
- decided-by [daemon-status-runs-cap](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-runs-cap.md)
- decided-by [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
