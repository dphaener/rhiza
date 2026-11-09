# request-slot-kept

capability · `01M6HT09X1W6W89R5N272GJXMP` · parent [plan-requests](../plan-requests.md)

## acceptance

daemon-request-slot.txtar with a fake rex: with max-runs 2, two ready change sets and no request, daemon --once starts both ship runs; with one queued request and max-runs 1 it starts one ship run and the request stays queued; with max-runs 2 it starts one ship run and the request's plan run, and daemon status shows the other change set as no-slot with 'waiting for a slot: kept for queued continuations or a plan request'; with max-runs 2, two ready change sets and only a queued milestone request, daemon --once starts the milestone request's plan run and one ship run (two runs, never three), and daemon status shows the other change set as no-slot; go test -trimpath ./cmd/rhiza -run TestScripts/daemon-request-slot passes

## statement

While a queued feature or bug plan request could start (no run, its retry time passed) and max-runs is above 1, the daemon's ship picks leave one slot for the request's plan run (besides any kept for queued continuations); a ready change set held back by this shows in daemon status as no-slot with the detail 'waiting for a slot: kept for queued continuations or a plan request'. A queued milestone request keeps no reserved slot: it takes the next free slot ahead of every other start. With max-runs 1, or no feature or bug request waiting, picks are unchanged

## Links

### Out

- decided-by [request-slot-kept-owner](../../../../../decision/request-slot-kept-owner.md)
- decided-by [milestone-interview-starts-at-once](../../../../../decision/milestone-interview-starts-at-once.md)
- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [pr-291-a-milestone-request](../../../../../evidence/pr-291-a-milestone-request.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6HSZMD2H2HZVGQQSVZZ4A1D` Plan requests are not starved by ready change sets: while a plan request waits, the daemon keeps one run slot for it
- `01M6K044DTSK7PT604NHPNGD17` Queued continuations are not starved by ready change sets: the daemon keeps a slot for each queued resume or retry
- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
