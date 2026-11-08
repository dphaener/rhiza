# m29-round-robin-slots

decision · `01M6RN37DJY5DVAKMTK7R1PWX0`

## question

When projects compete for the root daemon's run slots, how are slots shared?

## statement

Run slots are shared round-robin across projects under one machine-wide maxRuns, each project's next ready change set taken in its own roadmap order, with an optional per-project max.

## wording

Round-robin. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- relates-to [continuations-first-machine-wide](continuations-first-machine-wide.md) — Refines the owner-picked round-robin: rotation and roadmap order stay, but queued continuations, then queued starts and plan requests, go ahead of ready change sets machine-wide.
- decided-by [m29-one-daemon](../intent/milestone-29/m29-one-daemon.md)
- decided-by [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)
- decided-by [vm-second-project-ships](../capability/rhiza/orchestration/daemon/vm-second-project-ships.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
