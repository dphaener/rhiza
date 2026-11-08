# root-continuations-first

capability · `01M6XCJN76SFKMEST71CGEA6VG` · parent [daemon](../daemon.md)

## acceptance

internal/daemon roundrobin tests: with maxRuns 3, 2 running and one free slot, project A (turn next, 2 ready change sets) and project B (1 queued continuation), B's continuation takes the slot; with 2 free slots, A has a queued start and a ready change set and B has a continuation, the picks are B's continuation, then A's start; a project at its per-project max gets no continuation slot, and the slot goes to the next project's continuation, else its start; with no continuations queued, the existing roundrobin cases pick as before. daemon-round-robin.txtar and daemon-continue-slot.txtar still pass

## statement

The root daemon's round-robin (daemon.RoundRobin) fills freed machine-wide slots in three tiers: first every project's queued continuations that could start now, in round-robin order across projects; then queued starts and plan requests; then ready change sets. Each tier keeps each project's rhiza.project.maxRuns and the machine-wide maxRuns, and the rotation's last turn carries on as before. So a queued continuation never loses a freed slot to another project's start or roadmap pick. Counting is unchanged: an interrupted run takes no slot, and no run is pre-empted.

## Links

### Out

- decided-by [continuations-first-machine-wide](../../../../decision/continuations-first-machine-wide.md)

### In

- verifies [root-continuations-first-verified](../../../../evidence/root-continuations-first-verified.md)
- verifies [root-continuations-first-qa-1](../../../../evidence/root-continuations-first-qa-1.md)
- verifies [pr-393-fix-a-queued-continuation-can](../../../../evidence/pr-393-fix-a-queued-continuation-can.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
