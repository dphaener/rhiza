# request-slot-kept-owner

decision · `01M6HT0ZNHGD1V36Z7F77EXM52`

## statement

While a plan request waits, the daemon keeps one run slot for it: ship picks fill at most max-runs minus 1 slots. Owner, 2026-11-03: "I just filed one that I want moved to the top of the queue" and "yes, raise it to 3"; the request still got no slot because the new slot went to another ship run

## wording

the owner's words; the slot rule is the orchestrator's, as the owner's proxy

## Links

### In

- relates-to [continue-slot-kept-owner](continue-slot-kept-owner.md) — Generalises the plan request's kept slot to queued continuations
- decided-by [request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/request-slot-kept.md)

## History

- `01M6HSZMD2H2HZVGQQSVZZ4A1D` Plan requests are not starved by ready change sets: while a plan request waits, the daemon keeps one run slot for it
