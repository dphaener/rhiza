# requests-queue-local

decision · `01M6CS0MA6ZDN8MQ2GMYPA19ET`

## statement

A plan request is a line in a per-clone queue under .git/rhiza (requests.json), written by rhiza request <text> [--kind feature|bug|milestone]; the daemon, unless paused, takes the oldest queued request when a slot is free and starts a plan run with --new and the request's text (a bug's intent prefixed 'Fix:'), marking the request with its run; rhiza request ls shows each request's state (queued, planning in run R, planned into change set C, failed); the roadmap app's new-work keys file requests instead of opening a run

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m22-run-from-rex-ask](../source/m22-run-from-rex-ask.md)

### In

- relates-to [bug-intent-fix-no-colon](bug-intent-fix-no-colon.md) — changes requests-queue-local's 'Fix:' prefix to 'Fix '
- decided-by [plan-requests](../capability/rhiza/orchestration/daemon/plan-requests.md)
- decided-by [roadmap-new-work-requests](../capability/rhiza/orchestration/daemon/plan-requests/roadmap-new-work-requests.md)
- decided-by [requests-app](../capability/rhiza/orchestration/daemon/plan-requests/requests-app.md)
- decided-by [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6CRZP2E8H6EEEPVDAWAKQ6H` Plan milestone 22: run Rhiza from Rex, the owner's ask of 2026-11-01
