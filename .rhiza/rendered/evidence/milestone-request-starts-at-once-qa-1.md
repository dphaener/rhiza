# milestone-request-starts-at-once-qa-1

evidence · `01M6N3JEZDS3W3WHV34RR2AFDH` · actor [claude-code-23c10d47](../actor/claude-code-23c10d47.md)

## result

PASS

## statement

QA at dea4e5b6, built binary, fake rex/claude, real detached runs: with --max-runs 1 and a feature request's plan run already live, queuing a second feature then a milestone request, daemon --once started only the milestone request's plan run and request ls showed the second feature still queued; with one milestone interview live a second milestone request stayed queued after the next daemon --once.

## Links

### Out

- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
