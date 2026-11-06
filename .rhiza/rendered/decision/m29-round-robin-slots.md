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

- decided-by [m29-one-daemon](../intent/milestone-29/m29-one-daemon.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
