# projection-build-order

decision · `01M7GFT0CRDWZ6R37AKS2H5QPN`

## reasoning

Moving status, inbox, render, change ls and the daemon fold onto the projection is about six leaves of medium risk, and deleting the interim model touches the write path and every IndexPath branch; one PR holding both would be the largest of the milestone and hard to review or revert. Split, the first leaves the old path in place as a fallback while the reads move, and the second deletes it once nothing reads it. The parity harness is the projection's safety net, so the projection waits for it; check reads the whole log either way and the scripts change only the test harness, so both can build beside it. Notes keep the later work visible on the milestone without a draft that a backlog rule could pick up early.

## statement

Finishing the projection is two change sets: first every read moves onto one projection interface, then the interim model and the IndexPath branches are deleted. So milestone 39 builds five change sets: parity harness, check without rebuilds, reads from the projection, one read path, rhiza in-process in scripts. Reads from the projection depends on the parity harness, and one read path on reads from the projection; check and the in-process scripts depend on nothing. The two later items are notes on milestone 39, not draft change sets

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [milestone-39](../intent/milestone-39.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
