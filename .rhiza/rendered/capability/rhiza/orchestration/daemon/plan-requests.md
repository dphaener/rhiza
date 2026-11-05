# plan-requests

capability · `01M6CS230463DT24PWBJCPWD9Q` · parent [daemon](../daemon.md)

## acceptance

plan-requests.txtar with a fake rex: request queues an item shown queued by request ls; daemon --once starts run start --workflow plan --new with the request's text (a bug's prefixed "Fix ", so its title names the bug, a milestone's with --milestone), --autonomous only with plan auto and never for a milestone request, and request ls then shows planning with the run; when the run's journal names the planned change set, request ls shows planned with it; paused starts nothing; describe-golden.txtar is regenerated

## statement

rhiza request <text> [--kind feature|bug|milestone] queues a plan request in the clone, rhiza request ls shows each request's state, and the daemon, when not paused and a slot is free, starts a plan run for the oldest queued request, autonomous when plan autonomy is auto unless the request is of kind milestone, marking the request with its run and later the change set it planned

## Links

### Out

- decided-by [requests-queue-local](../../../../decision/requests-queue-local.md)

### In

- verifies [plan-requests-qa-2](../../../../evidence/plan-requests-qa-2.md)
- verifies [pr-191-plan-requests](../../../../evidence/pr-191-plan-requests.md)
- verifies [pr-289-milestone-planning-always](../../../../evidence/pr-289-milestone-planning-always.md)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
- `01M6GDJW3R8PXW7QBAZVH0QS4W` Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix'
- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
