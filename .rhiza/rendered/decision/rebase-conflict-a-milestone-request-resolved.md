# rebase-conflict-a-milestone-request-resolved

decision · `01M6NB42XNHYMYMK7P0S8T730W`

## statement

Resolve the a-milestone-request rebase conflicts by merging origin/main into cs/a-milestone-request: plan-requests.txtar keeps both sides (the milestone request starts first with --milestone and never --autonomous; the auto pass starts the feature, bug and search requests with --autonomous), plan-milestone-interview-waiting.txtar now expects run start --milestone --autonomous to be refused per milestone-planning-always (#289), and .rhiza/rendered is regenerated; merge commit 60a44d2c, scoped tests and rhiza check green

## wording

stated by the orchestrator acting as arbiter

## Links

### Out

- answers [rebase-conflict-a-milestone-request](../question/rebase-conflict-a-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
