# pr-289-milestone-planning-always

evidence · `01M6N6FT4WHQ1P0CWMPD0M14MX` · actor [dphaener](../actor/dphaener.md)

## head

e01b291bf168d3bda0ea60096d6782c9a2124a4a

## pr

https://github.com/dphaener/rhiza/pull/289

## result

1 check passed: check SUCCESS

## statement

Pull request #289 (https://github.com/dphaener/rhiza/pull/289) passed CI at e01b291

## status

verified

## Links

### Out

- verifies [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)
- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)
- verifies [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)
- verifies [plan-requests](../capability/rhiza/orchestration/daemon/plan-requests.md)
- verifies [plan-gate](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
