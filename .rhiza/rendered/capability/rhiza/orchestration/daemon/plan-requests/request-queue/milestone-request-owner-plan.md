# milestone-request-owner-plan

capability · `01M6MYG5CW1WV56MB997XMDE93` · parent [request-queue](../request-queue.md)

## acceptance

A testscript with a fake rex sets rhiza.autonomy.plan=auto and queues a milestone request and a feature request: the daemon's plan run start of the milestone request has --milestone and no --autonomous, its journal's started entry is not autonomous, and the feature request's run start carries --autonomous

## statement

The daemon starts the plan run of a milestone request without --autonomous even when rhiza.autonomy.plan is auto, so it waits on the owner's interview; feature and bug requests keep following rhiza.autonomy.plan

## Links

### Out

- decided-by [milestone-plan-never-autonomous](../../../../../../decision/milestone-plan-never-autonomous.md)
- derived-from [milestone-planning-always-plan](../../../../../../source/milestone-planning-always-plan.md)

### In

- verifies [milestone-request-owner-plan-verified](../../../../../../evidence/milestone-request-owner-plan-verified.md)
- verifies [milestone-request-owner-plan-qa-1](../../../../../../evidence/milestone-request-owner-plan-qa-1.md)
- verifies [milestone-request-owner-plan-qa-2](../../../../../../evidence/milestone-request-owner-plan-qa-2.md)
- verifies [pr-289-milestone-planning-always](../../../../../../evidence/pr-289-milestone-planning-always.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
