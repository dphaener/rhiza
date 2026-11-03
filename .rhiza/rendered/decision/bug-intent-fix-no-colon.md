# bug-intent-fix-no-colon

decision · `01M6GDTGZ5BC6QNNN01956K5T0`

## options

```json
[
  "Fix <text> (no colon)",
  "explicit title per request",
  "slug with a change-set ID suffix"
]
```

## statement

A bug request's plan run opens its change set with the intent "Fix <text>", not "Fix: <text>": the display title is the intent's first clause, so the colon made every bug's title "Fix", and two bug plan runs at once shared the branch cs/fix and worktree rhiza.fix (plan runs 01M6GCPN0H and 01M6GCQ7NJ collided on 2026-11-03; one failed at its gate)

## wording

the orchestrator's, deciding as the owner's proxy under the owner's standing instruction "act on my behalf" (2026-11-03); assumed, not the owner's words

## Links

### Out

- relates-to [requests-queue-local](requests-queue-local.md) — changes requests-queue-local's 'Fix:' prefix to 'Fix '

### In

- decided-by [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6GDJW3R8PXW7QBAZVH0QS4W` Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix'
