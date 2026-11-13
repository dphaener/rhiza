# spawn-starts-from-caller

decision · `01M7AHW6MHDSWVKY7PW6NR3PYT`

## statement

rhiza run spawn starts the child from its own process through the run's session-start path (as a leaves step's StartLeaf does) and journals the child into the run's journal without the run's lock, as rhiza run review-check does; the run's driver, on resume, records any child that ended while no spawn process watched it, once.

## wording

derived from the code by the planner of an autonomous plan run

## Links

### Out

- derived-from [plan-daemon-offers-spawn](../source/plan-daemon-offers-spawn.md)

### In

- decided-by [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)
- decided-by [spawn-durable](../capability/rhiza/orchestration/daemon/spawn-durable.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
