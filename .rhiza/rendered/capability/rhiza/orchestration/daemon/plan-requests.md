# plan-requests

capability · `01M6CS230463DT24PWBJCPWD9Q` · parent [daemon](../daemon.md)

## acceptance

plan-requests.txtar with a fake rex: request queues an item shown queued by request ls; daemon --once starts run start --workflow plan --new with the request's text (a bug's prefixed Fix:, a milestone's with --milestone), --autonomous only with plan auto, and request ls then shows planning with the run; when the run's journal names the planned change set, request ls shows planned with it; paused starts nothing; describe-golden.txtar is regenerated

## statement

rhiza request <text> [--kind feature|bug|milestone] queues a plan request in the clone, rhiza request ls shows each request's state, and the daemon, when not paused and a slot is free, starts a plan run for the oldest queued request, autonomous when plan autonomy is auto, marking the request with its run and later the change set it planned

## Links

### Out

- decided-by [requests-queue-local](../../../../decision/requests-queue-local.md)

### In

- verifies [plan-requests-qa-2](../../../../evidence/plan-requests-qa-2.md)
- verifies [pr-191-plan-requests](../../../../evidence/pr-191-plan-requests.md)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
