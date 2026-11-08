# request-queue

capability · `01M6D4SYDB1P7NX7DGKT0D19HM` · parent [plan-requests](../plan-requests.md)

## acceptance

plan-requests.txtar with a fake rex: request queues an item shown queued by request ls; daemon --once starts run start --workflow plan --new with the request's text (a bug's prefixed "Fix ", so its title names the bug, a milestone's with --milestone), --autonomous only with plan auto and never for a milestone request, and request ls then shows planning with the run; when the run's journal names the planned change set, request ls shows planned with it; paused starts nothing; describe-golden.txtar is regenerated

## statement

rhiza request <text> [--kind feature|bug|milestone] queues a plan request in the clone's .git/rhiza/requests.json, rhiza request ls shows each request's state, and the daemon, when not paused and a slot is free, starts a plan run for the oldest queued feature or bug request, autonomous when rhiza.autonomy.plan is auto, marking the request with its run and later the change set it planned; a queued milestone request instead takes the next free slot first, ahead of every other start and request, but waits while running is at rhiza.daemon.maxRuns (milestone-request-starts-at-once, plan-runs-share-max-runs), and its plan run always waits on the owner, never autonomous

## Links

### Out

- decided-by [requests-queue-local](../../../../../decision/requests-queue-local.md)
- decided-by [bug-intent-fix-no-colon](../../../../../decision/bug-intent-fix-no-colon.md)
- decided-by [milestone-interview-starts-at-once](../../../../../decision/milestone-interview-starts-at-once.md)
- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [request-queue-verified](../../../../../evidence/request-queue-verified.md)
- verifies [request-queue-qa-1](../../../../../evidence/request-queue-qa-1.md)
- verifies [pr-192-request-queue](../../../../../evidence/pr-192-request-queue.md)
- verifies [pr-289-milestone-planning-always](../../../../../evidence/pr-289-milestone-planning-always.md)
- verifies [pr-291-a-milestone-request](../../../../../evidence/pr-291-a-milestone-request.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../../evidence/pr-301-delete-the-old-rex-apps.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6D4SB0D8F9SPGK471DF9RV6` The request queue and the daemon turning requests into plan runs, as a leaf of plan requests
- `01M6GDJW3R8PXW7QBAZVH0QS4W` Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix'
- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
