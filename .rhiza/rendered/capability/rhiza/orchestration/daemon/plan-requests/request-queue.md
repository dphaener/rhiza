# request-queue

capability · `01M6D4SYDB1P7NX7DGKT0D19HM` · parent [plan-requests](../plan-requests.md)

## acceptance

plan-requests.txtar with a fake rex: request queues an item shown queued by request ls; daemon --once starts run start --workflow plan --new with the request's text (a bug's prefixed Fix:, a milestone's with --milestone), --autonomous only with plan auto, and request ls then shows planning with the run; when the run's journal names the planned change set, request ls shows planned with it; paused starts nothing; the roadmap app's n and m keys then reach a working request command; describe-golden.txtar is regenerated

## statement

rhiza request <text> [--kind feature|bug|milestone] queues a plan request in the clone's .git/rhiza/requests.json, rhiza request ls shows each request's state, and the daemon, when not paused and a slot is free, starts a plan run for the oldest queued request, autonomous when rhiza.autonomy.plan is auto, marking the request with its run and later the change set it planned

## Links

### Out

- decided-by [requests-queue-local](../../../../../decision/requests-queue-local.md)

### In

- verifies [request-queue-verified](../../../../../evidence/request-queue-verified.md)
- verifies [request-queue-qa-1](../../../../../evidence/request-queue-qa-1.md)
- verifies [pr-192-request-queue](../../../../../evidence/pr-192-request-queue.md)

## History

- `01M6D4SB0D8F9SPGK471DF9RV6` The request queue and the daemon turning requests into plan runs, as a leaf of plan requests
