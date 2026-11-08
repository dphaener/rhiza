# continuations-first-machine-wide

decision · `01M6XCGGE44R8KVNVM9PX33SDQ`

## reasoning

The question's own diagnosis, which the code confirms, is that the 75-minute wait came from three long runs holding maxRuns 3, and that daemonDemand counts a LiveInterrupted run as live but not as running, so the cap is not the bug. The real gap is in daemon.RoundRobin (internal/daemon/roundrobin.go): it takes each project's Ready list in project order, one item per project per round, so project B's queued 'start' or roadmap pick can take the slot that frees while project A's 'continue' waits. Each project already orders its own Ready list with continuations first (daemonDemand), and per-repo daemonPoll1 already keeps slots for queued continuations (daemonSlotKept), which shows the intent (#226/#243) that continuations go ahead of new starts. Making the machine-wide pass honour that same order is the smallest change that matches the intent. Pre-emption or reserved slots would go beyond the change set's intent and the owner's cap control (rhiza.daemon.maxRuns). The status change is cheap and fixes the second symptom the change set names: status listed the run only as 'interrupted', which hid that it was queued.

## statement

Fix two things and leave the cap alone. (1) Machine-wide priority: the root daemon's round-robin (daemon.RoundRobin) first hands freed slots to every project's queued continuations, in round-robin order across projects and within each project's rhiza.project.maxRuns, then to queued starts and plan requests, and only then to ready change sets; so a continuation never loses a freed slot to another project's plain start or roadmap pick. (2) Status: rhiza daemon status (and its JSON) shows an interrupted run with a queued continuation as queued, waiting for a slot, with its place in line, not only as 'interrupted'. Not in scope: counting rules (an interrupted run already takes no slot), pre-empting or pausing running runs, or a slot reserved outside maxRuns; a long wait while all slots hold long runs is expected behaviour.

## wording

(assumed, would have escalated)

## Links

### Out

- relates-to [m29-round-robin-slots](m29-round-robin-slots.md) — Refines the owner-picked round-robin: rotation and roadmap order stay, but queued continuations, then queued starts and plan requests, go ahead of ready change sets machine-wide.
- answers [question-29479c02-1](../question/question-29479c02-1.md)
- derived-from [queued-continuation-slot-interview](../source/queued-continuation-slot-interview.md)

### In

- decided-by [root-continuations-first](../capability/rhiza/orchestration/daemon/root-continuations-first.md)
- decided-by [daemon-status-queued-continuation](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-queued-continuation.md)
- decided-by [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
