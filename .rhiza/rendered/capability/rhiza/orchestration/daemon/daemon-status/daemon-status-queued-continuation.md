# daemon-status-queued-continuation

capability · `01M6XCJNBA3Y2XBG39EKDVCRZ7` · parent [daemon-status](../daemon-status.md)

## acceptance

daemon-continue-slot.txtar with a fake rex: with max-runs 1 and one running run, two interrupted runs with rhiza run resume queued show in rhiza daemon status as interrupted, queued, waiting for a slot, at places 1 and 2, and --json gives each a queued place of 1 and 2; an interrupted run with nothing queued shows no queued place; describe-golden.txtar regenerated for the new response field

## statement

rhiza daemon status, in text and in --json (each live run's entry, per project under the root daemon), shows an interrupted run that has a queued continuation that could start as queued, waiting for a slot, with its place in its project's continuation queue (1 for the oldest), as well as interrupted. An interrupted run with nothing queued shows as before.

## Links

### Out

- decided-by [continuations-first-machine-wide](../../../../../decision/continuations-first-machine-wide.md)

### In

- verifies [daemon-status-queued-continuation-verified](../../../../../evidence/daemon-status-queued-continuation-verified.md)
- verifies [daemon-status-queued-continuation-qa-1](../../../../../evidence/daemon-status-queued-continuation-qa-1.md)
- verifies [pr-393-fix-a-queued-continuation-can](../../../../../evidence/pr-393-fix-a-queued-continuation-can.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
