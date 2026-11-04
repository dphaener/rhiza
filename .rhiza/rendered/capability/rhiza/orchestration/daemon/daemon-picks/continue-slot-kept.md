# continue-slot-kept

capability · `01M6K04J5BCWYYP6PPKN7G267Y` · parent [daemon-picks](../daemon-picks.md)

## acceptance

daemon-continue-slot.txtar with a fake rex: with max-runs 2, two ready change sets and no continuation, daemon --once starts both ship runs; with a failed run's retry queued it starts one ship run and resumes the run with --retry, and daemon status shows the other change set as no-slot with 'waiting for a slot: kept for queued continuations or a plan request'; daemon-request-slot.txtar still passes with that wording

## statement

While queued continuations (rhiza run resume, approve, request-changes or resume --retry) could start now — their run exists, is not live, and has not ended unless the continuation is a retry, as the daemon continues them — the daemon's ship picks leave a slot for each, plus the one kept for a waiting plan request, never fewer than 0 picks; a ready change set held back by this shows in daemon status as no-slot with the detail 'waiting for a slot: kept for queued continuations or a plan request'. With nothing queued, picks are unchanged

## Links

### Out

- decided-by [continue-slot-kept-owner](../../../../../decision/continue-slot-kept-owner.md)

## History

- `01M6K044DTSK7PT604NHPNGD17` Queued continuations are not starved by ready change sets: the daemon keeps a slot for each queued resume or retry
