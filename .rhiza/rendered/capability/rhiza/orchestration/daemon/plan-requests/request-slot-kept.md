# request-slot-kept

capability · `01M6HT09X1W6W89R5N272GJXMP` · parent [plan-requests](../plan-requests.md)

## acceptance

daemon-request-slot.txtar with a fake rex: with max-runs 2, two ready change sets and no request, daemon --once starts both ship runs; with one queued request and max-runs 1 it starts one ship run and the request stays queued; with max-runs 2 it starts one ship run and the request's plan run, and daemon status shows the other change set as no-slot with 'waiting for a slot: kept for queued continuations or a plan request'

## statement

While a queued plan request could start (no run, its retry time passed) and max-runs is above 1, the daemon's ship picks leave one slot for the request's plan run (besides any kept for queued continuations); a ready change set held back by this shows in daemon status as no-slot with the detail 'waiting for a slot: kept for queued continuations or a plan request'. With max-runs 1, or no request waiting, picks are unchanged

## Links

### Out

- decided-by [request-slot-kept-owner](../../../../../decision/request-slot-kept-owner.md)

## History

- `01M6HSZMD2H2HZVGQQSVZZ4A1D` Plan requests are not starved by ready change sets: while a plan request waits, the daemon keeps one run slot for it
- `01M6K044DTSK7PT604NHPNGD17` Queued continuations are not starved by ready change sets: the daemon keeps a slot for each queued resume or retry
