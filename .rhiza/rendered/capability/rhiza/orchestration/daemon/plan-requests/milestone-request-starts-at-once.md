# milestone-request-starts-at-once

capability · `01M6MYSEB8N545M8M9APZTBXYS` · parent [plan-requests](../plan-requests.md)

## acceptance

cmd/rhiza/testdata/script/daemon-milestone-request.txtar (fake rex, rhiza.autonomy.plan owner) queues a feature then a milestone request: with --max-runs 1 and one run already live, daemon --once starts nothing and request ls shows both queued; with --max-runs 2 and one run live, daemon --once starts only the milestone request (run start --workflow plan --new ... --request <id> --milestone) and the feature stays queued; with one interrupted run and one free slot the milestone request starts and the run is not resumed; a second milestone request queued while the first's run is live is not started by the next --once; go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-milestone-request|plan-requests)' passes

## statement

The daemon starts the oldest queued milestone request's plan run in the next free run slot, ahead of every other start in the poll (resumes of interrupted runs included) and of every other queued request, but waits while the machine-wide slot count is at rhiza.daemon.maxRuns; once started it counts as a running run for later picks, and while one milestone request's plan run is live a second milestone request stays queued. A usage-reserve hold does not hold it; pause and low disk do

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../../decision/milestone-interview-starts-at-once.md)
- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [milestone-request-starts-at-once-verified](../../../../../evidence/milestone-request-starts-at-once-verified.md)
- verifies [milestone-request-starts-at-once-qa-1](../../../../../evidence/milestone-request-starts-at-once-qa-1.md)
- verifies [pr-291-a-milestone-request](../../../../../evidence/pr-291-a-milestone-request.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
