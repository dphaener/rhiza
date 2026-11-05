# milestone-request-starts-at-once

capability · `01M6MYSEB8N545M8M9APZTBXYS` · parent [plan-requests](../plan-requests.md)

## acceptance

A new testscript cmd/rhiza/testdata/script/daemon-milestone-request.txtar (fake rex, rhiza.autonomy.plan owner, --max-runs 1 with one run already live) queues a feature then a milestone request: daemon --once starts only the milestone request (run start --workflow plan --new ... --request <id> --milestone) and request ls shows the feature still queued; a second milestone request queued while the first's run is live is not started by the next --once; go test -trimpath ./cmd/rhiza -run TestScripts/daemon-milestone-request and TestScripts/plan-requests pass

## statement

The daemon starts the oldest queued milestone request's plan run at its next poll, ahead of every other queued request and whether or not a run slot is free; once started it counts as a running run for later picks, and while one milestone request's plan run is live a second milestone request stays queued

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../../decision/milestone-interview-starts-at-once.md)

### In

- verifies [milestone-request-starts-at-once-verified](../../../../../evidence/milestone-request-starts-at-once-verified.md)
- verifies [milestone-request-starts-at-once-qa-1](../../../../../evidence/milestone-request-starts-at-once-qa-1.md)
- verifies [pr-291-a-milestone-request](../../../../../evidence/pr-291-a-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
