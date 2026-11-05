# pr-291-a-milestone-request

evidence · `01M6NESMVYR3W7D00QD480TFWC` · actor [dphaener](../actor/dphaener.md)

## head

c059a7b848bf564eefcd824d8b931a5660428aac

## pr

https://github.com/dphaener/rhiza/pull/291

## result

1 check passed: check SUCCESS

## statement

Pull request #291 (https://github.com/dphaener/rhiza/pull/291) passed CI at c059a7b

## status

verified

## Links

### Out

- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- verifies [milestone-request-past-usage-hold](../capability/rhiza/orchestration/usage-governor/milestone-request-past-usage-hold.md)
- verifies [milestone-interview-waiting-at-launch](../capability/rhiza/orchestration/notifications/milestone-interview-waiting-at-launch.md)
- verifies [inbox-milestone-interview](../capability/rhiza/orchestration/inbox/inbox-milestone-interview.md)
- verifies [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)
- verifies [request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/request-slot-kept.md)
- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
