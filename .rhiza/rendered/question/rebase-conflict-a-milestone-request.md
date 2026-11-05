# rebase-conflict-a-milestone-request

question · `01M6N9WSQ04DB780K6M2QH4VAB`

## statement

Rebasing branch cs/a-milestone-request onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: cmd/rhiza/testdata/script/plan-requests.txtar?

## Links

### Out

- blocks change set `01M6MYFH4QYGSKDATBH509SEJF` (A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning")

### In

- answers [rebase-conflict-a-milestone-request-resolved](../decision/rebase-conflict-a-milestone-request-resolved.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
