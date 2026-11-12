# milestone-request-starts-at-once

capability · `01M6MYSEB8N545M8M9APZTBXYS` · parent [plan-requests](../plan-requests.md)

## acceptance

cmd/rhiza/testdata/script/daemon-milestone-request.txtar (fake rex, rhiza.autonomy.plan owner) queues a feature then a milestone request: with --max-runs 1 and one run already live, daemon --once starts nothing and request ls shows both queued; with --max-runs 2 and one run live, daemon --once starts only the milestone request (run start --workflow plan --new ... --request <id> --milestone) and the feature stays queued; with one interrupted run and one free slot the milestone request starts and the run is not resumed; a second milestone request queued while the first's plan run holds the owner (its start unjournaled within daemonStartGrace, or at its interview step) is not started by the next --once; go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-milestone-request|plan-requests)' passes

## statement

The daemon starts the oldest queued milestone request's plan run in the next free run slot, ahead of every other start in the poll (resumes of interrupted runs included) and of every other queued request, but waits while the machine-wide slot count is at rhiza.daemon.maxRuns; once started it counts as a running run for later picks, and a second milestone request stays queued while the first's plan run holds the owner (see milestone-wait-owner-steps). A usage-reserve hold does not hold it; pause and low disk do

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../../decision/milestone-interview-starts-at-once.md)
- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)
- decided-by [milestone-wait-owner-steps-only](../../../../../decision/milestone-wait-owner-steps-only.md)

### In

- verifies [milestone-request-starts-at-once-verified](../../../../../evidence/milestone-request-starts-at-once-verified.md)
- verifies [milestone-request-starts-at-once-qa-1](../../../../../evidence/milestone-request-starts-at-once-qa-1.md)
- verifies [pr-291-a-milestone-request](../../../../../evidence/pr-291-a-milestone-request.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)
- verifies [milestone-request-starts-at-once-qa-owner-steps-1](../../../../../evidence/milestone-request-starts-at-once-qa-owner-steps-1.md)
- verifies [pr-453-fix-a-queued-milestone-request](../../../../../evidence/pr-453-fix-a-queued-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
